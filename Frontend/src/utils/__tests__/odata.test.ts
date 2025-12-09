import { describe, it, expect } from 'vitest';
import {
  buildODataQuery,
  ODataFilter,
  odataEntityKey,
  odataAction,
  odataFunction,
  odataExpandWithOptions,
  parseODataCollection,
  parseODataEntity,
  type ODataCollectionResponse,
  type ODataSingleResponse,
} from '../odata';

describe('OData Utilities', () => {
  describe('buildODataQuery', () => {
    it('should build query with select', () => {
      const query = buildODataQuery({ select: ['Id', 'Name'] });
      expect(query).toBe('?$select=Id,Name');
    });

    it('should build query with expand', () => {
      const query = buildODataQuery({ expand: 'Club' });
      expect(query).toBe('?$expand=Club');
    });

    it('should build query with multiple expands', () => {
      const query = buildODataQuery({ expand: ['Club', 'User'] });
      expect(query).toBe('?$expand=Club,User');
    });

    it('should build query with filter', () => {
      const query = buildODataQuery({ filter: "Name eq 'test'" });
      expect(query).toBe("?$filter=Name eq 'test'");
    });

    it('should build query with orderby', () => {
      const query = buildODataQuery({ orderby: 'Name asc' });
      expect(query).toBe('?$orderby=Name asc');
    });

    it('should build query with pagination', () => {
      const query = buildODataQuery({ skip: 10, top: 5 });
      expect(query).toBe('?$skip=10&$top=5');
    });

    it('should build query with count', () => {
      const query = buildODataQuery({ count: true });
      expect(query).toBe('?$count=true');
    });

    it('should build query with search', () => {
      const query = buildODataQuery({ search: 'keyword' });
      expect(query).toBe('?$search=keyword');
    });

    it('should build complex query with multiple options', () => {
      const query = buildODataQuery({
        select: ['Id', 'Name'],
        expand: 'Club',
        filter: "Active eq true",
        orderby: 'CreatedAt desc',
        skip: 10,
        top: 20,
        count: true,
      });
      expect(query).toContain('$select=Id,Name');
      expect(query).toContain('$expand=Club');
      expect(query).toContain('$filter=Active eq true');
      expect(query).toContain('$orderby=CreatedAt desc');
      expect(query).toContain('$skip=10');
      expect(query).toContain('$top=20');
      expect(query).toContain('$count=true');
    });

    it('should return empty string for empty options', () => {
      const query = buildODataQuery({});
      expect(query).toBe('');
    });
  });

  describe('ODataFilter', () => {
    it('should create eq filter', () => {
      expect(ODataFilter.eq('Name', 'test')).toBe("Name eq 'test'");
      expect(ODataFilter.eq('Age', 25)).toBe('Age eq 25');
      expect(ODataFilter.eq('Active', true)).toBe('Active eq true');
    });

    it('should escape single quotes in strings', () => {
      expect(ODataFilter.eq('Name', "O'Brien")).toBe("Name eq 'O''Brien'");
    });

    it('should create ne filter', () => {
      expect(ODataFilter.ne('Status', 'deleted')).toBe("Status ne 'deleted'");
    });

    it('should create gt filter', () => {
      expect(ODataFilter.gt('Age', 18)).toBe('Age gt 18');
    });

    it('should create ge filter', () => {
      expect(ODataFilter.ge('Score', 100)).toBe('Score ge 100');
    });

    it('should create lt filter', () => {
      expect(ODataFilter.lt('Price', 50)).toBe('Price lt 50');
    });

    it('should create le filter', () => {
      expect(ODataFilter.le('Stock', 10)).toBe('Stock le 10');
    });

    it('should create and filter', () => {
      const filter = ODataFilter.and(
        ODataFilter.eq('Active', true),
        ODataFilter.gt('Age', 18)
      );
      expect(filter).toBe('(Active eq true) and (Age gt 18)');
    });

    it('should create or filter', () => {
      const filter = ODataFilter.or(
        ODataFilter.eq('Type', 'admin'),
        ODataFilter.eq('Type', 'owner')
      );
      expect(filter).toBe("(Type eq 'admin') or (Type eq 'owner')");
    });

    it('should create not filter', () => {
      const filter = ODataFilter.not(ODataFilter.eq('Deleted', true));
      expect(filter).toBe('not (Deleted eq true)');
    });

    it('should create contains filter', () => {
      expect(ODataFilter.contains('Name', 'test')).toBe("contains(Name, 'test')");
    });

    it('should create startswith filter', () => {
      expect(ODataFilter.startswith('Email', 'admin')).toBe("startswith(Email, 'admin')");
    });

    it('should create endswith filter', () => {
      expect(ODataFilter.endswith('Email', '@example.com')).toBe("endswith(Email, '@example.com')");
    });

    it('should create isNull filter', () => {
      expect(ODataFilter.isNull('DeletedAt')).toBe('DeletedAt eq null');
    });

    it('should create isNotNull filter', () => {
      expect(ODataFilter.isNotNull('CreatedAt')).toBe('CreatedAt ne null');
    });

    it('should create in filter', () => {
      expect(ODataFilter.in('Status', ['active', 'pending'])).toBe("Status in ('active', 'pending')");
      expect(ODataFilter.in('Age', [18, 21, 25])).toBe('Age in (18, 21, 25)');
    });

    it('should create complex nested filters', () => {
      const filter = ODataFilter.and(
        ODataFilter.or(
          ODataFilter.eq('Type', 'admin'),
          ODataFilter.eq('Type', 'owner')
        ),
        ODataFilter.gt('Age', 18),
        ODataFilter.isNotNull('Email')
      );
      expect(filter).toBe("((Type eq 'admin') or (Type eq 'owner')) and (Age gt 18) and (Email ne null)");
    });
  });

  describe('odataEntityKey', () => {
    it('should create entity key reference', () => {
      const key = odataEntityKey('Clubs', 'abc-123');
      expect(key).toBe("Clubs('abc-123')");
    });
  });

  describe('odataAction', () => {
    it('should create action URL', () => {
      const url = odataAction('Clubs', 'abc-123', 'Leave');
      expect(url).toBe("Clubs('abc-123')/Leave");
    });
  });

  describe('odataFunction', () => {
    it('should create bound function URL without params', () => {
      const url = odataFunction('Clubs', 'abc-123', 'IsAdmin');
      expect(url).toBe("Clubs('abc-123')/IsAdmin");
    });

    it('should create bound function URL with params', () => {
      const url = odataFunction('Events', 'event-1', 'ExpandRecurrence', {
        startDate: '2024-01-01',
        endDate: '2024-12-31',
      });
      expect(url).toBe("Events('event-1')/ExpandRecurrence(startDate='2024-01-01',endDate='2024-12-31')");
    });

    it('should create unbound function URL', () => {
      const url = odataFunction(null, null, 'GetDashboardNews');
      expect(url).toBe('GetDashboardNews');
    });

    it('should create unbound function URL with params', () => {
      const url = odataFunction(null, null, 'SearchGlobal', { query: 'test' });
      expect(url).toBe("SearchGlobal(query='test')");
    });
  });

  describe('odataExpandWithOptions', () => {
    it('should create simple expand', () => {
      const expand = odataExpandWithOptions('Members');
      expect(expand).toBe('Members');
    });

    it('should create expand with filter', () => {
      const expand = odataExpandWithOptions('Members', {
        filter: "Role eq 'admin'",
      });
      expect(expand).toBe("Members($filter=Role eq 'admin')");
    });

    it('should create expand with select', () => {
      const expand = odataExpandWithOptions('Members', {
        select: ['Id', 'Name'],
      });
      expect(expand).toBe('Members($select=Id,Name)');
    });

    it('should create expand with multiple options', () => {
      const expand = odataExpandWithOptions('Events', {
        filter: 'StartTime gt 2024-01-01',
        orderby: 'StartTime asc',
        top: 10,
      });
      expect(expand).toContain('Events(');
      expect(expand).toContain('$filter=');
      expect(expand).toContain('$orderby=');
      expect(expand).toContain('$top=');
    });
  });

  describe('parseODataCollection', () => {
    it('should extract value array from OData response', () => {
      const response: ODataCollectionResponse<{ id: string; name: string }> = {
        '@odata.context': 'http://example.com/$metadata#Clubs',
        '@odata.count': 2,
        value: [
          { id: '1', name: 'Club 1' },
          { id: '2', name: 'Club 2' },
        ],
      };
      
      const result = parseODataCollection(response);
      expect(result).toEqual([
        { id: '1', name: 'Club 1' },
        { id: '2', name: 'Club 2' },
      ]);
    });
  });

  describe('parseODataEntity', () => {
    it('should extract entity and remove metadata', () => {
      const response: ODataSingleResponse<{ id: string; name: string }> = {
        '@odata.context': 'http://example.com/$metadata#Clubs/$entity',
        id: '1',
        name: 'Test Club',
      };
      
      const result = parseODataEntity(response);
      expect(result).toEqual({
        id: '1',
        name: 'Test Club',
      });
      expect('@odata.context' in result).toBe(false);
    });
  });
});
