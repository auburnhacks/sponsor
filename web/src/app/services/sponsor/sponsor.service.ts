import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { AuthService } from '../auth/auth.service';

@Injectable({
  providedIn: 'root'
})
export class SponsorService {

  constructor(
    private http: HttpClient,
    private authService: AuthService
  ) { }

  getCompanies(): Observable<Company[]> {
    return this.http.get<Company[]>(environment.apiBase + "/sponsor/companies", this.getHttpOptions())
  }

  createSponsor(formData): Observable<Sponsor> {
    const sponsorData = {
      sponsor: {
        name: formData.sponsorName,
        email: formData.sponsorEmail,
        password: formData.sponsorPassword,
        ACL: formData.aclListStr,
        company: {
          id: formData.companyId
        }
      }
    };
    return this.http.post<Sponsor>(environment.apiBase + "/sponsor", sponsorData, this.getHttpOptions());
  }

  private getHttpOptions() {
    return {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + this.authService.user().token,
      })
    }
  }
}
