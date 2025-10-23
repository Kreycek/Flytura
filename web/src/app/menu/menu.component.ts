import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-menu',
  standalone: true,
  imports: [CommonModule,RouterModule,TranslateModule],
  templateUrl: './menu.component.html',
  styleUrl: './menu.component.css'
})
export class MenuComponent {
  i=false;

  @Input() activeMenu: boolean = false;

  constructor(
    private router: Router,
  ) {}
  ngOnInit() {

    setTimeout(() => {
      this.i=true;
    }, 1000);
  }

  selectedMenu=''
  navigateToPages(url:string,menuId: string) {
    this.selectedMenu = menuId;
    this.router.navigate([url]);
  }

}
