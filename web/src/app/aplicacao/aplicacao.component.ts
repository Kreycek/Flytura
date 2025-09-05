import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { MenuComponent } from '../menu/menu.component';
import { CenterComponent } from '../center/center.component';
import { FooterComponent } from '../footer/footer.component';
import { NavbarComponent } from '../navbar/navbar.component';

@Component({
  selector: 'app-aplicacao',
  standalone: true,
  imports: [CommonModule, RouterOutlet, MenuComponent, FooterComponent, NavbarComponent],
  templateUrl: './aplicacao.component.html',
  styleUrl: './aplicacao.component.css'
})
export class AplicacaoComponent {

}
