import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth/auth.service';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { User } from '../models/user.model';
import { FormGroup } from '@angular/forms';

@Component({
  selector: 'app-admin',
  templateUrl: './admin.component.html',
  styleUrls: ['./admin.component.css']
})
export class AdminComponent implements OnInit {

  public user: User;
  public addSponsorForm: FormGroup;
  
  constructor(private authService: AuthService, private activatedRouter: ActivatedRoute,
              private router: Router) { }

  ngOnInit() {
    this.activatedRouter.params.subscribe((params: Params): void => {
      if(params['id'] == undefined) {
        this.router.navigate(['/login']);
      }
      this.user = this.authService.user();
    });
  }

}
