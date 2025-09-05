import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';

@Component({
  selector: 'app-menu',
  standalone: true,
  imports: [CommonModule,RouterModule],
  templateUrl: './menu.component.html',
  styleUrl: './menu.component.css'
})
export class MenuComponent {
  i=false;

  constructor(private router: Router) {}
  ngOnInit() {
    setTimeout(() => {
      this.i=true;
    }, 1000);
  }

  navigateToPages(url:string) {
    this.router.navigate([url]);
  }

}
