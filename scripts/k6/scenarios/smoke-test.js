import http from 'k6/http';
import { sleep } from 'k6';

// export const options = {
//   // A number specifying the number of VUs to run concurrently.
//   vus: 10,
//   // A string specifying the total duration of the test run.
//   duration: '30s',
// };

// export default function () {
//   http.get('http://localhost:7777');
//   sleep(1);
// }

export const options = {
  vus: 5,
  duration: '1m',
  thresholds: {
    'http_req_duration{status:200}': ['max>=0'],
    'http_req_duration{status:429}': ['max>=0'],
    'http_req_duration{status:500}': ['max>=0'],
  },
  'summaryTrendStats': ['min', 'med', 'avg', 'p(90)', 'p(95)', 'max', 'count'],
};

export default function () {
  http.get('http://rate-limiter-api:8080');

  http.get('http://rate-limiter-api:8080', {
    headers: {
      'API_KEY': 'abc123'
    }
  });
}


