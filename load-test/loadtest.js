import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 1000 },
        { duration: '1m', target: 1000 }, 
        { duration: '30s', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<50'],
        http_req_failed: ['rate<0.0001'],
    },
};

const BASE_URL = 'http://localhost:8080';
const USERNAME = 'testuser';
const PASSWORD = 'testpassword';

function authenticate() {
    const url = `${BASE_URL}/api/auth`;
    const payload = JSON.stringify({
        username: USERNAME,
        password: PASSWORD,
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(url, payload, params);
    check(res, {
        'auth status is 200': (r) => r.status === 200,
    });

    return res.json().token;
}

export default function () {
    const token = authenticate();

    const infoUrl = `${BASE_URL}/api/info`;
    const infoParams = {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    };
    const infoRes = http.get(infoUrl, infoParams);
    check(infoRes, {
        'info status is 200': (r) => r.status === 200,
    });

    const sendCoinUrl = `${BASE_URL}/api/sendCoin`;
    const sendCoinPayload = JSON.stringify({
        toUser: 'anotheruser',
        amount: 10,
    });
    const sendCoinParams = {
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
        },
    };
    const sendCoinRes = http.post(sendCoinUrl, sendCoinPayload, sendCoinParams);
    check(sendCoinRes, {
        'sendCoin status is 200': (r) => r.status === 200,
    });

    sleep(1);
}