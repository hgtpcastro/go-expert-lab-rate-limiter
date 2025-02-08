import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 200 },
    { duration: '3m', target: 200 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    'http_req_duration{status:200}': ['max>=0'],
    'http_req_duration{status:429}': ['max>=0'],
    'http_req_duration{status:500}': ['max>=0'],
  },
  'summaryTrendStats': ['min', 'med', 'avg', 'p(90)', 'p(95)', 'max', 'count'],
};

export default function () {
  http.get('http://localhost:8080');

  http.get('http://localhost:8080', {
    headers: {
      'API_KEY': 'any-api-key'
    }
  });

  sleep(1);
}