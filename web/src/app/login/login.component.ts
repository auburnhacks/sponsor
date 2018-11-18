import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AuthService } from '../services/auth/auth.service';
import { Admin, Sponsor } from '../models/user.model';
import { Error } from '../models/error.model';
import { Router, ActivatedRoute } from '@angular/router';
import { filter } from 'rxjs/operators';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  public loginError: string = "";

  public loginOptions = ['sponsor', 'admin'];

  public loginForm: FormGroup;

  constructor(private fb: FormBuilder, private authService: AuthService, private router: Router,
              private acRoute: ActivatedRoute) { 
    this.loginForm = this.fb.group({
      rememberMe: [false],
      source: [this.loginOptions[0]],
      email: ['', Validators.required],
      password: ['', Validators.required],
    });
  }

  ngOnInit() {
    this.acRoute.queryParams
        .pipe(filter(params => params.action))
        .subscribe(params => {
          console.log(params)
          if (params['action'] == "logout") {
            this.authService.logout();
            return;
          }
        });
    if (this.authService.isAuthenticated()) {
      this.router.navigate(['/home', this.authService.user().id]);
    }
  }

  login() {
    const formValues = this.loginForm.value;
    if (formValues['source'] == "admin") {
      this.authService
      .loginAdmin(formValues['email'], formValues['password'])
      .then((user: Admin) => {
          this.router.navigate(['/home', user.id]);
        },
        (reason: Error) => {
          this.loginError = reason.message;
      });
    } else if (formValues['source'] == "sponsor") {
      this.authService
      .loginSponsor(formValues['email'], formValues['password'])
      .then((user: Sponsor) => {
        console.log(user);
        this.router.navigate(['/home', user.id]);
      },
      (reason: Error) => {
        this.loginError = reason.message;
      });
    }
  }

}
