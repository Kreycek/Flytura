import { Component } from '@angular/core';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { InvoicesService } from '../../modulos/gestion/invoices/invoices.service';
import { react } from 'plotly.js';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-total-group',
  imports: [TranslateModule, CommonModule],
  templateUrl: './total-group.component.html',
  styleUrl: './total-group.component.css'
})
export class TotalGroupComponent {
  constructor(
private translate: TranslateService,
private invoicesService:InvoicesService

  ) {

  }

  TotalLast30Days=0
  TotalGeral=0
  dataRangeValues:any[]=[]

  ngOnInit() {
    this.invoicesService.getCountLast30DaysByDtImports().subscribe((res=>{
     

      this.TotalLast30Days=res.totalLast30Days

    }))


     this.invoicesService.getCountLast30DaysByAmountRange().subscribe((res=>{
      console.log('teste ',res);

      this.TotalGeral=res.totalGeral
      this.dataRangeValues=res.data
     

    }))
  } 
}
