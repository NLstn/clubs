/**
 * OData v4 Query Builder Utilities
 * 
 * Helper functions for constructing OData v4 query URLs with proper encoding
 * and type safety. These utilities help migrate from REST API v1 to OData v2.
 * 
 * @see https://www.odata.org/documentation/
 */

/**
 * OData query options interface
 */
export interface ODataQueryOptions {
  /** Select specific fields: $select=Field1,Field2 */
  select?: string[];
  /** Expand related entities: $expand=RelatedEntity */
  expand?: string | string[];
  /** Filter results: $filter=Field eq 'value' */
  filter?: string;
  /** Order by fields: $orderby=Field asc */
  orderby?: string;
  /** Skip N results for pagination: $skip=10 */
  skip?: number;
  /** Take N results: $top=10 */
  top?: number;
  /** Include total count: $count=true */
  count?: boolean;
  /** Search across entity: $search="search term" */
  search?: string;
}

/**
 * Builds an OData query string from options
 * 
 * @param options - OData query options
 * @returns Query string (e.g., "?$filter=Name eq 'test'&$top=10")
 */
export function buildODataQuery(options: ODataQueryOptions): string {
  const parts: string[] = [];

  if (options.select && options.select.length > 0) {
    parts.push(`$select=${options.select.join(',')}`);
  }

  if (options.expand) {
    const expandValue = Array.isArray(options.expand) 
      ? options.expand.join(',') 
      : options.expand;
    parts.push(`$expand=${expandValue}`);
  }

  if (options.filter) {
    parts.push(`$filter=${options.filter}`);
  }

  if (options.orderby) {
    parts.push(`$orderby=${options.orderby}`);
  }

  if (options.skip !== undefined) {
    parts.push(`$skip=${options.skip}`);
  }

  if (options.top !== undefined) {
    parts.push(`$top=${options.top}`);
  }

  if (options.count) {
    parts.push('$count=true');
  }

  if (options.search) {
    parts.push(`$search=${options.search}`);
  }

  return parts.length > 0 ? `?${parts.join('&')}` : '';
}

/**
 * Filter expression builder helpers
 */
export const ODataFilter = {
  /** Field equals value: Field eq 'value' */
  eq: (field: string, value: string | number | boolean) => 
    `${field} eq ${formatValue(value)}`,
  
  /** Field not equals value: Field ne 'value' */
  ne: (field: string, value: string | number | boolean) => 
    `${field} ne ${formatValue(value)}`,
  
  /** Field greater than value: Field gt 5 */
  gt: (field: string, value: string | number) => 
    `${field} gt ${formatValue(value)}`,
  
  /** Field greater than or equal: Field ge 5 */
  ge: (field: string, value: string | number) => 
    `${field} ge ${formatValue(value)}`,
  
  /** Field less than value: Field lt 5 */
  lt: (field: string, value: string | number) => 
    `${field} lt ${formatValue(value)}`,
  
  /** Field less than or equal: Field le 5 */
  le: (field: string, value: string | number) => 
    `${field} le ${formatValue(value)}`,
  
  /** Logical AND: (expr1) and (expr2) */
  and: (...expressions: string[]) => 
    expressions.map(e => `(${e})`).join(' and '),
  
  /** Logical OR: (expr1) or (expr2) */
  or: (...expressions: string[]) => 
    expressions.map(e => `(${e})`).join(' or '),
  
  /** Logical NOT: not (expr) */
  not: (expression: string) => 
    `not (${expression})`,
  
  /** String contains: contains(Field, 'value') */
  contains: (field: string, value: string) => 
    `contains(${field}, ${formatValue(value)})`,
  
  /** String starts with: startswith(Field, 'value') */
  startswith: (field: string, value: string) => 
    `startswith(${field}, ${formatValue(value)})`,
  
  /** String ends with: endswith(Field, 'value') */
  endswith: (field: string, value: string) => 
    `endswith(${field}, ${formatValue(value)})`,
  
  /** Field is null: Field eq null */
  isNull: (field: string) => 
    `${field} eq null`,
  
  /** Field is not null: Field ne null */
  isNotNull: (field: string) => 
    `${field} ne null`,
  
  /** Field in list: Field in ('val1', 'val2') */
  in: (field: string, values: (string | number)[]) => 
    `${field} in (${values.map(formatValue).join(', ')})`,
};

/**
 * Format value for OData query (adds quotes for strings)
 */
function formatValue(value: string | number | boolean): string {
  if (typeof value === 'string') {
    // Escape single quotes in strings
    const escaped = value.replace(/'/g, "''");
    return `'${escaped}'`;
  }
  return value.toString();
}

/**
 * Build OData entity key: EntitySet('key-value')
 * 
 * @param entitySet - The entity set name (e.g., "Clubs")
 * @param key - The entity key value
 * @returns OData entity reference (e.g., "Clubs('abc-123')")
 */
export function odataEntityKey(entitySet: string, key: string): string {
  return `${entitySet}('${key}')`;
}

/**
 * Build OData action call: EntitySet('key')/ActionName
 * 
 * @param entitySet - The entity set name
 * @param key - The entity key value
 * @param action - The action name
 * @returns OData action URL (e.g., "Clubs('abc-123')/Leave")
 */
export function odataAction(entitySet: string, key: string, action: string): string {
  return `${odataEntityKey(entitySet, key)}/${action}`;
}

/**
 * Build OData function call: EntitySet('key')/FunctionName
 * 
 * @param entitySet - The entity set name
 * @param key - The entity key value (optional for unbound functions)
 * @param functionName - The function name
 * @param params - Function parameters as key-value pairs
 * @returns OData function URL
 */
export function odataFunction(
  entitySet: string | null,
  key: string | null,
  functionName: string,
  params?: Record<string, string | number | boolean>
): string {
  let url = entitySet && key 
    ? `${odataEntityKey(entitySet, key)}/${functionName}` 
    : functionName;
  
  if (params && Object.keys(params).length > 0) {
    const paramString = Object.entries(params)
      .map(([k, v]) => `${k}=${formatValue(v)}`)
      .join(',');
    url += `(${paramString})`;
  }
  
  return url;
}

/**
 * OData expand with nested query options
 * Example: $expand=Members($filter=Role eq 'admin';$select=Id,Name)
 */
export function odataExpandWithOptions(
  entity: string,
  options?: ODataQueryOptions
): string {
  if (!options) {
    return entity;
  }
  
  const nestedQuery = buildODataQuery(options).replace('?', '');
  return nestedQuery ? `${entity}(${nestedQuery})` : entity;
}

/**
 * Common OData response interfaces
 */
export interface ODataCollectionResponse<T> {
  '@odata.context'?: string;
  '@odata.count'?: number;
  '@odata.nextLink'?: string;
  value: T[];
}

export type ODataSingleResponse<T> = T & {
  '@odata.context'?: string;
};

/**
 * Parse OData collection response and return just the value array
 */
export function parseODataCollection<T>(response: ODataCollectionResponse<T>): T[] {
  return response.value;
}

/**
 * Extract entity from OData single entity response
 */
export function parseODataEntity<T>(response: ODataSingleResponse<T>): T {
  // Remove OData metadata fields
  const { '@odata.context': _unused, ...entity } = response;
  void _unused; // Mark as intentionally unused
  return entity as T;
}
