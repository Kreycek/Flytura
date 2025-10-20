import { Component, EventEmitter, Output } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { jwtDecode } from 'jwt-decode';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [RouterModule,TranslateModule],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.css'
})
export class NavbarComponent {
  @Output() am = new EventEmitter<boolean>
  private isActive = true;

  welcomeUser=''


  constructor(private translate: TranslateService,) {
    
interface MeuPayload {

  name: string;
  lastName: string;
  
  // outros campos que seu token tiver
}

const token = localStorage.getItem('token'); // ou onde você armazenou o JWT

    if (token) {
      const decoded = jwtDecode<MeuPayload>(token);
      console.log(decoded);

      this.welcomeUser=' ' + decoded.name + ' ' + decoded.lastName
    }
    
  }

  changeLanguage(language:string) {
    localStorage.removeItem('language')
    localStorage.setItem('language', language);
     this.translate.use(language);
        
  }

  activeMenu() {
      this.isActive = !this.isActive;
    this.am.emit(this.isActive);
  }
}
