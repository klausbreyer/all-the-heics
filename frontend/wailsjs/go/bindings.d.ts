import * as models from './models';

export interface go {
  "main": {
    "App": {
		List():Promise<Array<models.heic>>
		WorkFile(arg1:string):Promise<boolean>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
