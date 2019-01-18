import { Component, OnInit, AfterContentInit } from '@angular/core';
import { Router, ActivatedRoute, Params } from '@angular/router';
import { AuthService } from '../services/auth/auth.service';
import { User } from '../models/user.model';
import { ParticipantService } from '../services/participant/participant.service';
import { Participant } from '../models/participant.model';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit, AfterContentInit {

  public id: string;
  public user: User;
  public participants: Array<Participant>;

  
  constructor(private acRoute: ActivatedRoute, public authService: AuthService,
              private router: Router, private participantService: ParticipantService) { }

  public ngOnInit(): void {
    this.acRoute.params.subscribe((params) : void => {
      // Check if the current user id is authenticated
      if (params['id']) {
        if (!this.authService.isAuthenticated()) {
          this.router.navigate(['/login']);
        }
        // Since the user might come back after closing the browser
        // validating them to make sure they are authorized to look at the homepage
        this.id = params['id'];
        this.authService
          .validateUser(this.id)
          .then((isValid) => {
            if(!isValid) {
              console.info('Invalid userId: ' + this.id);
              this.router.navigate(['/login'])
            }
            this.user = this.authService.user();
          })
      } else {
        this.router.navigate(['/login']);
      }
    });
  }

  public ngAfterContentInit(): void {
    this.participantService
        .list()
        .then((particiIn: Array<Participant>) => {
          this.participants = particiIn.filter(p => p.name && p.resume);
        }, 
        (reason) => {
          console.log(reason);
        });
  }
  
  public hasUser() {
    return this.user !== undefined;
  }

  public getUserId() {
    if(this.user) {
      return this.user.id;
    }
    return undefined;
  }
}
