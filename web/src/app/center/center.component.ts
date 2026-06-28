import { Component, ViewChild } from '@angular/core';
import { ProcessedComponent } from './processed/processed.component';
import { CommonModule } from '@angular/common';
import { GraficImportAndStatusComponent } from "./grafic-import-and-status/grafic-import-and-status.component";
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { TotalGroupComponent } from "./total-group/total-group.component";

@Component({
  selector: 'app-center',
  standalone: true,
  imports: [ProcessedComponent, CommonModule, GraficImportAndStatusComponent, TranslateModule, TotalGroupComponent],
  
  providers: [CommonModule],
  templateUrl: './center.component.html',
  styleUrl: './center.component.css'
})
export class CenterComponent {

     constructor( 
       
     
         private translate: TranslateService
      ) {}

  showProcessed = false;
  showGraficDateImport=false
  totals=true;
  

    toggleProcessed() {
      this.showProcessed = !this.showProcessed;
      this.showGraficDateImport =false;    
      this.totals=false
    }

    toggleImportaStatus() {    
      this.showGraficDateImport = !this.showGraficDateImport;
      this.showProcessed = false;
      this.totals=false
    }

     toggleImportaTotal() {    
      this.showGraficDateImport = false;
      this.showProcessed = false;
      this.totals=true
    }

}
