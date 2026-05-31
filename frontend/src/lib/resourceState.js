export const RESOURCE_STATUS = {
  idle: 'idle',
  loading: 'loading',
  success: 'success',
  empty: 'empty',
  error: 'error',
}

export function createResourceState(status = RESOURCE_STATUS.idle, payload = {}) {
  return {
    status,
    data: null,
    error: null,
    ...payload,
  }
}
