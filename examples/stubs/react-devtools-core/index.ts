/**
 * Stand-in for Ink's optional React DevTools integration. Without it,
 * `bun build --compile` hoists the import and every binary fails at startup.
 * Ink only calls these when DEV=true and a DevTools server is reachable.
 */
const devtools = {
  initialize(): void {},
  connectToDevTools(): void {},
};

export default devtools;
