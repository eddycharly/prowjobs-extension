const defaultOptions = {
  method: 'get'
};

export function getHeaders(headers = {}) {
  return {
    Accept: 'application/json',
    'Content-Type': 'application/json',
    ...headers
  };
}

export function checkStatus(response = {}) {
  if (response.ok) {
    switch (response.status) {
      case 201: 
        return response.headers;
      case 204: 
        return {};
      default: {
        let responseAsJson = response.json();
        return responseAsJson;
      }
    }
  }

  const error = new Error(response.statusText);
  error.response = response;
  throw error;
}

export function request(uri, options = defaultOptions) {
  return fetch(uri, {
    ...options
  }).then(checkStatus);
}

export function get(uri) {
  return request(uri, {
    method: 'get',
    headers: getHeaders()
  });
}

export function post(uri, body) {
  return request(uri, {
    method: 'post',
    headers: getHeaders(),
    body: JSON.stringify(body)
  });
}

export function put(uri, body) {
  return request(uri, {
    method: 'put',
    headers: getHeaders(),
    body: JSON.stringify(body)
  });
}

export function deleteRequest(uri) {
  return request(uri, {
    method: 'delete',
    headers: getHeaders()
  });
}
