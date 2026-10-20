import http from 'k6/http'
import { check, sleep } from 'k6';

const BASE_URL = __ENV.BASE_URL;
const jsonHeaders = { headers: { 'Content-Type': 'application/json' } };

export const options = {
    stages: [
        { duration: '1m', target: 10},
        { duration: '3m', target: 50},
        { duration: '1m', target: 0},
    ],
    thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: ['p(95)<500', 'p(99)<1000'],
        checks: ['rate>0.99'],
    },
};

export default function() {
    const suffix = `${Date.now()}-${__VU}-${__ITER}`;
    const email = `k6-${suffix}@example.com`;
    const password = 'K6Password123!';

    const register = http.post(
        `${BASE_URL}/api/auth/register`,
        JSON.stringify({
            email,
            password,
        }),
        jsonHeaders,
    );

    check(register, {
        'register returns 201': (response) => response.status === 201,
    });

    const login = http.post(
        `${BASE_URL}/api/auth/login`,
        JSON.stringify({
            email,
            password,
        }),
        jsonHeaders,
    );

    check(login, {
        'login returns 200': (response) => response.status === 200,
        'login returns a token': (response) => Boolean(response.json('token')),
    })

    const token = login.json('token');

    if (!token) {
        return;
    }

    const authHeaders = {
        headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
        },
    };

    const create = http.post(
        `${BASE_URL}/api/recipes`,
        JSON.stringify({
        title: `k6 recipe ${suffix}`,
        description: 'Recipe created during k6 performance testing',
        ingredients: 'Flour, water, salt',
        }),
        authHeaders,
    );

    check(create, {
        'create recipe returns 201': (response) => response.status === 201,
    });

    const recipeId = create.json('id');
    if (!recipeId) {
        return;
    }

    const list = http.get(`${BASE_URL}/api/recipes`);
    check(list, {
        'list recipes returns 200': (response) => response.status === 200,
    });

    const search = http.get(
        `${BASE_URL}/api/recipes?search=${encodeURIComponent('k6 recipe')}`,
    );
    check(search, {
        'search recipes returns 200': (response) => response.status === 200,
    });

    const recipe = http.get(`${BASE_URL}/api/recipes/${recipeId}`);
    check(recipe, {
        'get recipe returns 200': (response) => response.status === 200,
    });

    const remove = http.del(`${BASE_URL}/api/recipes/${recipeId}`, null, authHeaders);
    check(remove, {
        'delete recipe returns 204': (response) => response.status === 204,
    });

    sleep(1);
}