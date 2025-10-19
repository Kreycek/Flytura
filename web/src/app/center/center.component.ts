import { Component, ViewChild } from '@angular/core';

import { BrowserModule } from '@angular/platform-browser';
// @ts-ignore
import * as Plotly from 'plotly.js-dist-min';
import { InvoicesService } from '../modulos/companys/invoices/invoices.service';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { AirLineService } from '../modulos/companys/airLine/airLIne.service';
import { ConfigService } from '../services/config.service';
import { MAT_DATE_LOCALE, MatNativeDateModule } from '@angular/material/core';
import moment from 'moment';
import { ModalOkComponent } from '../modal/modal-ok/modal-ok.component';
import { TranslateModule } from '@ngx-translate/core';
@Component({
  selector: 'app-center',
  standalone: true,
  imports: [CommonModule,FormsModule,MatDatepickerModule,TranslateModule,MatNativeDateModule],
  
  providers: [{ provide: MAT_DATE_LOCALE, useValue: 'pt-BR' }],
  templateUrl: './center.component.html',
  styleUrl: './center.component.css'
})
export class CenterComponent {
@ViewChild(ModalOkComponent) modalOk!: ModalOkComponent;  

    searchAirlineCode:string='';
    searchAirlineDtInicio:string='';
    searchAirlineDtFim :string='';
    searchStatus:string='';
    statusImportData:any[]=[];
    airLInes:any[]=[];
    msgNotFound=false;

    constructor( 
        private invoiceService: InvoicesService,
        private airLineService: AirLineService,
        public configService:ConfigService) {}

      async invalidDate(date: string, msg:string, showAlert:boolean=false) : Promise<boolean> {              
          const value = moment(date); // inputDate pode ser string, Date, etc.
  
          if (!value.isValid()) {
            if(showAlert) {
                const resultado = await this.modalOk.openModal(msg,true);             
                    if (resultado) {
                  
                      // Insira aqui a lógica para continuar após a confirmação
                    } else {
                      
                    }
                }
                  
            return true;
          } else {
           return false;
          }
      }

      graphicData(_values:any[], _labels:any[]) {
           const data = [{
            values: _values,
            labels: _labels,
            type: 'pie',
            // textinfo: 'label+value+percent',
            textinfo: 'label',
            insidetextorientation: 'radial'
          }];

          const layout = {
            title: 'Gráfico em Pizza com Legenda à Direita',
            height: 300,
            legend: {
              orientation: 'v',
              x: 0.62,
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
      }

  async generateGraphic() {

          let _values:any[]=[]
          let _labels:any[]=[]
          let objPesquisar: {                    
                companyCode?:string| null;  
                startDate?:string | null;
                endDate?:string | null;
                status?:string | null;              
            }
      
            if(this.searchAirlineDtInicio || this.searchAirlineDtFim ){
           
                if(await this.invalidDate(this.searchAirlineDtInicio, 'Data de início inválida, ou se a data de fim estiver preenchida é necessário preencher uma data de início.', true)) {
                    return false;
                }
      
                if(await this.invalidDate(this.searchAirlineDtFim, 'Data de fim inválida, ou se a data de início estiver preenchida é necessário preencher uma data de fim.', true)) {
                    return false;
                }
                    
                const startDate = moment(this.searchAirlineDtInicio, 'YYYY-MM-DD');
                const endDate = moment(this.searchAirlineDtFim, 'YYYY-MM-DD');
      
                if (startDate.isAfter(endDate)) {
                    const resultado = await this.modalOk.openModal('Data de início não pode ser maior que a data de fim',true);             
                        if (resultado) {
                          return false;
                          // Insira aqui a lógica para continuar após a confirmação
                        } else {
                          
                        }
                  
                }
      
                if (endDate.isBefore(startDate)) {
                  const resultado = await this.modalOk.openModal('Data de fim é menor que a data de inicio',true);             
                        if (resultado) {
                          return false;
                          // Insira aqui a lógica para continuar após a confirmação
                        } else {
                          
                        }          
                }
            }

             objPesquisar= { 
                  
                    companyCode:this.searchAirlineCode,
                    startDate:this.searchAirlineDtInicio ? moment(new Date(this.searchAirlineDtInicio)).format('YYYY-MM-DDT00:00:00Z') : null,
                    endDate:this.searchAirlineDtFim ? moment(new Date(this.searchAirlineDtFim)).format('YYYY-MM-DDT23:59:59Z') :null,
                    status:this.searchStatus
                  };

        this.invoiceService.GroupByCompanyName(objPesquisar.status,objPesquisar.companyCode,objPesquisar.startDate,objPesquisar.endDate).subscribe((response:any)=>{

         if(response) {
          response.forEach((element:any) => {
              _values.push(element.total)
              _labels.push(element._id + ' ' + element.total)
          });

          
          this.graphicData(_values,_labels);
          this.msgNotFound=false;
        }
        else {
          this.msgNotFound=true;
          this.graphicData([],[]);
        }

       

          
        
      })

      return null
  }

  ngOnInit() {

      let _values:any[]=[]
      let _labels:any[]=[]

      this.invoiceService.GroupByCompanyName().subscribe((response:any)=>{

       
          response.forEach((element:any) => {
              _values.push(element.total)
              _labels.push(element._id + ' ' + element.total)
          });

          if(_values.length>0 && _labels.length>0)
                this.graphicData(_values,_labels);

          
      })

        this.airLineService.getAllAirLine().subscribe((response)=>{
          this.airLInes=response;

        });

        this.invoiceService.getAllStatusImportData().subscribe((response:any)=>{     
        this.statusImportData=response

          });
  }
}
