import http from 'k6/http';
import {sleep, check} from 'k6';

export const options = {
    vus: 200,
    stages: [
        // need to add ramp-up and ramp-down stages
        {duration: '30s', target: 100},
    ],
};

export default function () {
    const payload = JSON.stringify({
        user_id: '408f6222-c435-4535-9de8-bf4ca22a79bc',
        item_id: 'f7e61eef-6f7e-4f50-9f75-19c64b06a7f2',
        amount: 1,
    });
    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };
    let res = http.post('http://localhost:3000/orders', payload, params);
    check(res, {"status is 200": (res) => res.status === 200});
    //sleep(1);
}
