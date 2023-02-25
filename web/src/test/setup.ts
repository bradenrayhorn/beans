import "@testing-library/jest-dom";
import { setupServer } from "msw/node";
import "cross-fetch/polyfill";

export const server = setupServer();

export const api = (path: string) => `http://127.0.0.1:8000${path}`;

beforeAll(() => server.listen({ onUnhandledRequest: "error" }));
afterAll(() => server.close());
afterEach(() => server.resetHandlers());
