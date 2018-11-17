import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

import { LoginService } from '../services/login/login.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  public loginError: string;

  public loginOptions = ['sponsor', 'admin'];

  public loginForm: FormGroup;

  constructor(private fb: FormBuilder, private loginService: LoginService) { 
    this.loginForm = this.fb.group({
      rememberMe: [false],
      source: [this.loginOptions[0]],
      email: ['', Validators.required],
      password: ['', Validators.required],
    });
  }

  ngOnInit() {
  }

  login() {
    console.info(this.loginForm.value);
  }

}
