import * as models from './models';

export interface go {
  "main": {
    "App": {
		Greet():Promise<string>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
