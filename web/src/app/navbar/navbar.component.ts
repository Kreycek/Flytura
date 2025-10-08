import { Component, EventEmitter, Output } from '@angular/core';
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
  @Output() am = new EventEmitter<boolean>
  private isActive = true;


  constructor(private translate: TranslateService) {
 
    
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
