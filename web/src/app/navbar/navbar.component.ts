import { AfterViewInit, Component, EventEmitter, Output, ViewChild } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { jwtDecode } from 'jwt-decode';
import { ConfigService } from '../services/config.service';
import { ModalOkComponent } from '../modal/modal-ok/modal-ok.component';
import packageJson from '../../../package.json';

// ...

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [RouterModule, TranslateModule, ModalOkComponent],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.css'
})
export class NavbarComponent  {
  @ViewChild(ModalOkComponent) modalOk!: ModalOkComponent;  
  @Output() am = new EventEmitter<boolean>();
  private isActive = true;

  welcomeUser = '';
  version=''

  constructor(
    private translate: TranslateService, 
    private router: Router,
     private configService:ConfigService
  ) {}

   ngOnInit() {
     this.version = packageJson.version;
    
     const tokenData=this.configService.returnTokenData()
     this.welcomeUser=' ' + tokenData?.name + ' ' + tokenData?.lastName
  }



  changeLanguage(language: string) {
    localStorage.removeItem('language');
    localStorage.setItem('language', language);
    this.translate.use(language);
  }

  async activeMenu() {
    this.isActive = !this.isActive;
    this.am.emit(this.isActive);
 
  }

  navigateToPages(url: string) {   
    this.router.navigate([url]);
  }

  logout() {
    localStorage.removeItem('token');
    window.location.assign('https://flytura.com');
  }
}
