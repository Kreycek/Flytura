import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, HttpClientModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = 'Flytura';

      constructor(private translate: TranslateService, ) {
        
      }

      ngOnInit() {
        this.setLanguage();
      }

      setLanguage() {

        let language:string | null = localStorage.getItem('language');
        if(language) {  
            this.translate.use(language); // ou 'en', 'es', etc.
        }
        else {
            localStorage.setItem('language', 'es');
            this.translate.use('es'); 
        }
      }


}
