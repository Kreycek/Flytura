import { Component, ViewChild } from '@angular/core';
import { ProcessedComponent } from './processed/processed.component';
import { CommonModule } from '@angular/common';
import { GraficImportAndStatusComponent } from "./grafic-import-and-status/grafic-import-and-status.component";
import { TranslateModule, TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-center',
  standalone: true,
  imports: [ProcessedComponent, CommonModule, GraficImportAndStatusComponent,TranslateModule],
  
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

    toggleProcessed() {
      this.showProcessed = !this.showProcessed;
      this.showGraficDateImport =false;    
    }

    toggleImportaStatus() {    
      this.showGraficDateImport = !this.showGraficDateImport;
      this.showProcessed = false;
    }

}
