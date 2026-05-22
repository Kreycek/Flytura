import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TranslateModule } from '@ngx-translate/core';
@Component({
  selector: 'app-alert-more-columns',
  imports: [TranslateModule],
  templateUrl: './alert-more-columns.component.html',
  styleUrl: './alert-more-columns.component.css'
})
export class AlertMoreColumnsComponent {
      @Input() viewAlert : boolean = false;

          
  ngOnInit() {
    
      if (window.innerWidth < 5000) {
          this.viewAlert = true;

          setTimeout(() => {
            this.viewAlert = false;
          }, 4000);
        }
      }
}
