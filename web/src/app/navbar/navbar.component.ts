import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [RouterModule],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.css'
})
export class NavbarComponent {

  /**
   *
   */
  constructor(private translate: TranslateService) {
 
    
  }

  changeLanguage(language:string) {
    localStorage.removeItem('language')
    localStorage.setItem('language', language);
     this.translate.use(language);
        
  }
}
