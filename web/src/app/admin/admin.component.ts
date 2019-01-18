import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth/auth.service';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { User } from '../models/user.model';
import { FormGroup, FormBuilder, Validators, FormControl, FormArray } from '@angular/forms';
import { SponsorService } from '../services/sponsor/sponsor.service';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-admin',
  templateUrl: './admin.component.html',
  styleUrls: ['./admin.component.css']
})
export class AdminComponent implements OnInit {

  public user: User;
  public aclListMap = [
    { id: 'read', name: 'Read' },
    { id: 'write', name: 'Write' },
    { id: 'download', name: 'Download' },
  ]
  public addSponsorForm: FormGroup;
  public companies: Company[] = new Array<Company>();
  public showAddAlert: boolean = false;

  constructor(
    public authService: AuthService,
    private activatedRouter: ActivatedRoute,
    private router: Router,
    private fb: FormBuilder,
    private sponsorService: SponsorService) { 

    let aclControls = this.aclListMap.map(c => new FormControl(false));
    aclControls[0].setValue(true); // set read to always be true
    this.addSponsorForm = this.fb.group({
      companyId: [],
      aclListMap: new FormArray(aclControls),
      sponsorName: ['', Validators.required],
      sponsorEmail: ['', Validators.required],
      sponsorPassword: ['', Validators.required],
    });
  }

  ngOnInit() {
    this.activatedRouter.params.subscribe((params: Params): void => {
      if(params['id'] == undefined) {
        this.router.navigate(['/login']);
      }
      this.user = this.authService.user();
      // Load the companies as soon as the component mounts on screen
      this.sponsorService.getCompanies().subscribe((data) => {
        if (!data['companies']) {
          return;
        }
        data['companies'].forEach((c: Company) => {
          return this.companies.push(c);
        });
      },
      (error: any) => console.log(error),
      () => { 
        this.addSponsorForm.setControl('companyId', new FormControl(this.companies[0].id));
        this.addSponsorForm.controls['companyId'].setValue(this.companies[0].id, {onlySelf: true})
      });
    });
  }



  createSponsor()  {
    console.log(this.addSponsorForm.value);
    if (!this.addSponsorForm.valid) {
      // TODO: change this to a error component that is mounted everytime
      // an error occurs
      console.log('form not valid');
    } else {
      // Change the aclListMap from the selected true and false to the actual
      // list
      console.log('We can submit this sponsor');
      this.boolsToACLStr(this.addSponsorForm.value['aclListMap'])
        .toPromise()
        .then((aclStr: string) => {
          // Update the sponsor form with the actual string
          this.addSponsorForm.setControl('aclListStr', new FormControl(aclStr));
          console.log(this.addSponsorForm.value);
          this.sponsorService.createSponsor(this.addSponsorForm.value)
            .subscribe((sp: Sponsor) => {
              // updating UI based on the result
              this.showAddAlert = true;
              // console.log(sp);
            });
        });
    }
  }

  generateUniquePassword() {
    let randomPassword = Math.random().toString(36).slice(-8);
    this.addSponsorForm.controls['sponsorPassword'].setValue(randomPassword);
  }

  boolsToACLStr(boolACLList: boolean[]): Observable<string> {
    let aclStr = new Observable<string>((observer) => {
      let acls = [];
      boolACLList.forEach((isActive, i) => {
        if(isActive) {
          acls.push(this.aclListMap[i].id);
        }
      });
      observer.next(acls.join(","));
      observer.complete();
    });
    return aclStr;
  }
  
  get aclFormData() {
    return this.addSponsorForm.controls.aclListMap;
  }

  getUserId() {
    return this.user.id;
  }

  showAddSponsorAlert() {
    return this.showAddAlert;
  }
}
