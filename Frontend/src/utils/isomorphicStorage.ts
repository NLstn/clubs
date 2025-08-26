const isBrowser = typeof window !== 'undefined';

export const getItem = (key: string): string | null => {
  if (!isBrowser) return null;
  return window.localStorage.getItem(key);
};

export const setItem = (key: string, value: string): void => {
  if (!isBrowser) return;
  window.localStorage.setItem(key, value);
};

export const removeItem = (key: string): void => {
  if (!isBrowser) return;
  window.localStorage.removeItem(key);
};

export default { getItem, setItem, removeItem };
