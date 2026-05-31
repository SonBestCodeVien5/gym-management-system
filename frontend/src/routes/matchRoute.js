import { APP_HOME_PATH, APP_PATH, APP_ROUTES, LOGIN_PATH, NOT_FOUND_ROUTE, ROOT_PATH } from './routeConfig.js'

export function normalizePath(pathname) {
  if (!pathname || pathname === ROOT_PATH) {
    return ROOT_PATH
  }

  const normalized = pathname.replace(/\/+$/, '')
  return normalized || ROOT_PATH
}

function decodeParamSegment(segment) {
  try {
    return decodeURIComponent(segment)
  } catch {
    return null
  }
}

function matchPattern(pattern, pathname) {
  const patternParts = normalizePath(pattern).split('/').filter(Boolean)
  const pathParts = normalizePath(pathname).split('/').filter(Boolean)

  if (patternParts.length !== pathParts.length) {
    return null
  }

  const params = {}

  for (let index = 0; index < patternParts.length; index += 1) {
    const patternPart = patternParts[index]
    const pathPart = pathParts[index]

    if (patternPart.startsWith(':')) {
      const decodedParam = decodeParamSegment(pathPart)

      if (decodedParam === null) {
        return null
      }

      params[patternPart.slice(1)] = decodedParam
      continue
    }

    if (patternPart !== pathPart) {
      return null
    }
  }

  return params
}

export function matchRoute(pathname) {
  const path = normalizePath(pathname)

  if (path === ROOT_PATH) {
    return { type: 'root', path, route: null, params: {} }
  }

  if (path === LOGIN_PATH) {
    return { type: 'login', path, route: null, params: {} }
  }

  if (path === APP_PATH) {
    return { type: 'redirect', path, route: null, params: {}, redirectTo: APP_HOME_PATH }
  }

  for (const route of APP_ROUTES) {
    const params = matchPattern(route.path, path)

    if (params) {
      return { type: 'app', path, route, params }
    }
  }

  if (path.startsWith(`${APP_PATH}/`)) {
    return { type: 'app-not-found', path, route: NOT_FOUND_ROUTE, params: {} }
  }

  return { type: 'redirect', path, route: null, params: {}, redirectTo: ROOT_PATH }
}
