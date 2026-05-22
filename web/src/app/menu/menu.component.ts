import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { jwtDecode } from 'jwt-decode';
import { EnumPerfil } from '../modulos/Enum/perfil';

@Component({
  selector: 'app-menu',
  standalone: true,
  imports: [CommonModule,RouterModule,TranslateModule],
  templateUrl: './menu.component.html',
  styleUrl: './menu.component.css'
})
export class MenuComponent {
 

  @Input() activeMenu: boolean = false;
  adm:boolean=true

  constructor(
    private router: Router,
  ) {}
  ngOnInit() {

    const token = localStorage.getItem('token'); // ou onde você armazenou o JWT
          
              if (token) {
                const decoded = jwtDecode<any>(token);
                if(decoded && decoded.perfis && decoded.perfis.length>0) {
                    this.adm=decoded.perfis.some((data:any)=> data===EnumPerfil.ADM)
                }
              }


  }

  selectedMenu=''
  navigateToPages(url:string,menuId: string) {
    this.selectedMenu = menuId;
    this.router.navigate([url]);
  }

}
