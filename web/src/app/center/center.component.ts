import { Component, ViewChild } from '@angular/core';

import { BrowserModule } from '@angular/platform-browser';
// @ts-ignore
import * as Plotly from 'plotly.js-dist-min';
import { InvoicesService } from '../modulos/companys/invoices/invoices.service';

@Component({
  selector: 'app-center',
  standalone: true,
  imports: [],
  templateUrl: './center.component.html',
  styleUrl: './center.component.css'
})
export class CenterComponent {

  /**
   *
   */
  constructor( private invoiceService: InvoicesService,) {
    
    
  }

  ngOnInit() {

      let _values:any[]=[]
      let _labels:any[]=[]

      this.invoiceService.GroupByCompanyName().subscribe((response:any)=>{

          console.log('response ',response);
          response.forEach((element:any) => {
              _values.push(element.total)
              _labels.push(element._id + ' ' + element.total)
          });

          
          const data = [{
            values: _values,
            labels: _labels,
            type: 'pie',
            textinfo: 'label+value+percent',
            insidetextorientation: 'radial'
          }];

          const layout = {
            title: 'Gráfico em Pizza com Legenda à Direita',
            height: 300,
            legend: {
              orientation: 'v',
              x: 0.60,
              y: 0.5,
              xanchor: 'left'
            },
            margin: {
              l: 50,
              r: 10,
              t: 50,
              b: 50
            }
          };

          Plotly.newPlot('graficoPlotly', data, layout);   
      })
  }
}
