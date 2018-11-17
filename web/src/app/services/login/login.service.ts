import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class LoginService {

  constructor() { }

  public login(email: string, password: string): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      
    });
  }
}
