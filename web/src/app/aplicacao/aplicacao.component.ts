import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { MenuComponent } from '../menu/menu.component';
import { CenterComponent } from '../center/center.component';
import { FooterComponent } from '../footer/footer.component';
import { NavbarComponent } from '../navbar/navbar.component';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-aplicacao',
  standalone: true,
  imports: [CommonModule, RouterOutlet, MenuComponent, FooterComponent, NavbarComponent],
  templateUrl: './aplicacao.component.html',
  styleUrl: './aplicacao.component.css'
})
export class AplicacaoComponent {
  constructor( private translate: TranslateService) {

  }

  am:boolean=true
  ngOnInit() {
console.log('am',this.am);
    let language:string | null = localStorage.getItem('language');
    if(language) {  
       this.translate.use(language); // ou 'en', 'es', etc.
       }
       else {
        localStorage.setItem('language', 'es');
        this.translate.use('es'); 
       }
  }

  activeMenu(value:boolean) {
    this.am=value;
    console.log('click ',this.am);
  }
}
