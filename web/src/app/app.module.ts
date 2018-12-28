import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RouterModule, Routes } from '@angular/router';
import { ClarityModule, ClrPasswordContainer, ClrFormsNextModule } from '@clr/angular';

import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { LoginComponent } from './login/login.component';
import { AuthService } from './services/auth/auth.service';
import { AdminComponent } from './admin/admin.component';
import { SponsorService } from './services/sponsor/sponsor.service';
import { UpdateProfileComponent } from './update-profile/update-profile.component';

const appRoutes: Routes = [
  { path:'', component: HomeComponent },
  { path: 'login', component: LoginComponent },
  { path: 'home/:id', component: HomeComponent },
  { path: 'logout', redirectTo: '/login?action=logout'},
  { path: 'admin/:id', component: AdminComponent },
  { path: 'profile/:id/update', component: UpdateProfileComponent},
]

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    LoginComponent,
    AdminComponent,
    UpdateProfileComponent
  ],
  imports: [
    FormsModule,
    ClarityModule,
    ClrFormsNextModule,
    BrowserModule,
    HttpClientModule,
    ReactiveFormsModule,
    BrowserAnimationsModule,
    RouterModule.forRoot(appRoutes),
  ],
  providers: [
    AuthService,
    SponsorService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
