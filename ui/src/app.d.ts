import "unplugin-icons/types/svelte";
// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
  namespace App {
    // interface Error {}
    interface Locals {
      isLoggedIn: boolean;
      maybeUser: User | undefiend;
    }
    // interface PageData {}
    // interface Platform {}
  }
}

type User = {
  id: string;
  username: string;
};

export { User };
