import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth/auth.service';
import { User } from '../models/user.model';
import { FormGroup, FormBuilder, FormsModule } from '@angular/forms';

@Component({
  selector: 'app-update-profile',
  templateUrl: './update-profile.component.html',
  styleUrls: ['./update-profile.component.css']
})
export class UpdateProfileComponent implements OnInit {

  public user: User;
  public updateUserForm: FormGroup;
  
  constructor(public authService: AuthService, private fb: FormBuilder) { }

  ngOnInit() {
    this.user = this.authService.user();
    this.updateUserForm = this.fb.group({
      name: [this.user.name],
      email: [this.user.email],
      password: [],
      confirmPassword: [],
    });
    // users cannot change their email address for now
    this.updateUserForm.controls['email'].disable();
  }

  updateProfile() {
    let newPassword = this.updateUserForm.value['password'];
    let newConfirmPassword = this.updateUserForm.value['confirmPassword'];
    if (newPassword !== newConfirmPassword) {
      console.log('do elegant error handling here');
      return;
    }
    // update user data here
  }
}
