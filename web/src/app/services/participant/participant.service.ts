import { Injectable } from '@angular/core';
import { Participant } from '../../models/participant.model';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { AuthService } from '../auth/auth.service';

@Injectable({
  providedIn: 'root'
})
export class ParticipantService {

  constructor(private http: HttpClient, private authService: AuthService) { }

  public list(): Promise<Array<Participant>> {
    return new Promise<Array<Participant>>((resolve, reject) => {
      this.http.get(environment.apiBase + "/sponsor/participants", 
          { headers: new HttpHeaders().append("Authorization", "Bearer " + this.authService.user().token)})
          .toPromise()
          .then(
            (data) => {
              let participants = new Array<Participant>();
              for (let participant of data['participants']) {
                participants.push(participant as Participant);
              }
              resolve(participants);
            },
            (reason) => reject(reason.error as Error))
    });
  }
}
