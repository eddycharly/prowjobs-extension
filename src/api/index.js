import { get, deleteRequest } from './comms';

const dashboardAPIRoot = getDashboardAPIRoot();

export function getDashboardAPIRoot() {
  const { href, hash } = window.location;
  let baseURL = href.replace(hash, '');
  if (baseURL.endsWith('/')) {
    baseURL = baseURL.slice(0, -1);
  }
  return baseURL;
}

export function getProwjobs() {
  const uri = `${dashboardAPIRoot}/proxy/apis/prow.k8s.io/v1/prowjobs`;
  return get(uri);
}

export function deleteProwjob(prowjob) {
  const uri = `${dashboardAPIRoot}/proxy/apis/prow.k8s.io/v1/namespaces/${prowjob.metadata.namespace}/prowjobs/${prowjob.metadata.name}`;
  return deleteRequest(uri);
}
