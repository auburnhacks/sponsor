import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, Params } from '@angular/router';
import { AuthService } from '../services/auth/auth.service';
import { User } from '../models/user.model';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  public id: string;
  public user: User;

  
  constructor(private acRoute: ActivatedRoute, private authService: AuthService,
              private router: Router) { }

  public ngOnInit(): void {
    this.acRoute.params.subscribe((params) : void => {
      // Check if the current user id is authenticated
      if (params['id']) {
        if (!this.authService.isAuthenticated()) {
          this.router.navigate(['/login']);
        }
        this.id = params['id'];
        this.user = this.authService.user();
      } else {
        this.router.navigate(['/login']);
      }
    });
  }

}
