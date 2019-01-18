import { Injectable } from '@angular/core';
import { Admin, Sponsor, User } from '../../models/user.model';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Error } from '../../models/error.model';
import * as moment from "moment";
import { Type } from '@angular/compiler';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  private currentUser: Admin | Sponsor;
  private type: string;

  constructor(private http: HttpClient) { 
    if (this.isAuthenticated() && !this.currentUser) {
      this.loadCurrentUser();
    }
  }

  public user(): User {
    return this.currentUser;
  }

  public isAdmin(): boolean {
    return this.type == "Admin" ? true : false;
  }

  public logout() {
    localStorage.clear();  
  }

  public isAuthenticated(): boolean {
    if (moment().isSameOrBefore(this.getTokenExpiry()) && this.getToken().length > 0) {
      return true;
    }
    return false;
  }

  public loginAdmin(email: string, password: string): Promise<Admin> {
    return new Promise<Admin>((resolve, reject) => {
      this.http.post(environment.apiBase + "/sponsor/admin/login", {email, password}, 
        { headers: new HttpHeaders().append("Content-Type", "application/json")})
        .toPromise()
        .then((data) => {
          let admin = data['admin'] as Admin;
          admin.token = data['token'];
          this.setUser(admin, "Admin");
          this.setSession(admin.token);
          resolve(admin);
        },
        (reason) => {
          reject(reason.error as Error);
        });
    });
  }

  public loginSponsor(email: string, password: string): Promise<Sponsor> {
    return new Promise<Sponsor>((resolve, reject) => {
      this.http.post(environment.apiBase + "/sponsor/login", 
        {email, "password_plain_text": password}, 
        { headers: new HttpHeaders().append("Content-Type", "application/json")})
        .toPromise()
        .then((data) => {
          let sp = data['sponsor'] as Sponsor;
          sp.token = data['token'];
          this.setUser(sp, "Sponsor");
          this.setSession(sp.token);
          resolve(sp);
        },
        (reason) => {
          reject(reason.error as Error);
        });
    });
  }

  /**
   * validateUser sends a request to the server to check and see if the given token is still
   * valid
   * @param userId 
   */
  public validateUser(userId: string): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      if (userId.length == 0) {
        reject(false);
      }
      if (this.isAdmin()) {
        this.http
        .get(environment.apiBase + "/sponsor/admin/" + userId,
          { headers: new HttpHeaders().append("Authorization", "Bearer " + this.getToken())})
        .toPromise()
        .then((adminData) => {
          resolve(true);
        },
        (reason) => reject(false));
      } else {
        // check sponsor endpoint as loggedin user is not a admin
        this.http
          .get(environment.apiBase + "/sponsor/" + userId + "/info",
          { headers: new HttpHeaders().append("Authorization", "Bearer " + this.getToken()) })
          .toPromise()
          .then((sponsorData) => {
            resolve(true);
          }, (reason) => reject(false));
      }
    });
  }

  public updateUser(changeData): Observable<Admin | Sponsor> {
    let updatedUser = new Observable<Admin | Sponsor>((observer) => {
      
    });
    return updatedUser;
  }

  private loadCurrentUser() {
    const savedUser = JSON.parse(localStorage.getItem("sess_user"));
    if (savedUser.type == "Admin") {
      this.currentUser = savedUser.obj as Admin;
      this.type = "Admin";
      return;
    }
    this.currentUser = savedUser.obj as Sponsor;
    this.type = "Sponsor";
  }

  private setUser(user: Admin | Sponsor, type: string): boolean {
    const savedUser = { type: type, obj: user};
    localStorage.setItem("sess_user", JSON.stringify(savedUser));
    this.loadCurrentUser();
    return true;
  }

  private setSession(token: string) {
    const expiresIn = moment().add(1, "hour").toISOString();
    localStorage.setItem("token", token);
    localStorage.setItem("token_expires_at", expiresIn);
  }

  private getTokenExpiry() {
    const expiration = localStorage.getItem("token_expires_at");
    return moment(expiration);
  }

  private getToken(): string {
    return localStorage.getItem("token");
  }
}
