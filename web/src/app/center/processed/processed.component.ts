import { Component, ViewChild } from '@angular/core';
// @ts-ignore
import * as Plotly from 'plotly.js-dist-min';
import { PurcharseRecordService } from '../../modulos/Account/purcharseRecord/purcharse-record.service';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { AirLineService } from '../../modulos/gestion/airLine/airLIne.service';
import { ConfigService } from '../../services/config.service';
import { MAT_DATE_LOCALE, MatNativeDateModule } from '@angular/material/core';
import moment from 'moment';
import { ModalOkComponent } from '../../modal/modal-ok/modal-ok.component';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ModuloService } from '../../modulos/modulo.service';

@Component({
  selector: 'app-processed',
  standalone: true,
  providers: [{ provide: MAT_DATE_LOCALE, useValue: 'pt-BR' }],
  imports: [CommonModule, FormsModule, MatDatepickerModule, TranslateModule, MatNativeDateModule, ModalOkComponent],
  templateUrl: './processed.component.html',
  styleUrl: './processed.component.css'
})
export class ProcessedComponent {
@ViewChild(ModalOkComponent) modalOk!: ModalOkComponent;  

    searchAirlineCode:string='';
    searchAirlineDtInicio:string='';
    searchAirlineDtFim :string='';
    searchStatus:string='';
    statusImportData:any[]=[];
    airLInes:any[]=[];
    msgNotFound=false;

    constructor( 
        private purcharseRecordService: PurcharseRecordService,
        private airLineService: AirLineService,
        public configService:ConfigService,
        public moduloService:ModuloService,
         private translate: TranslateService
      ) {}

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

     
 graphicData(_values: any[], _labels: any[], graphTarget: string | HTMLElement) {
    // sanity check de dados
    if (!Array.isArray(_values) || !Array.isArray(_labels) || _values.length !== _labels.length || _values.length === 0) {
      console.warn('Dados inválidos para Plotly', { _values, _labels });
      return;
    }

    // sanity check do target
    const el = typeof graphTarget === 'string'
      ? document.getElementById(graphTarget)
      : graphTarget;

    if (!el) {
      console.warn('Elemento do gráfico não encontrado:', graphTarget);
      return;
    }

    const data = [{
      values: _values,
      labels: _labels,
      type: 'pie',
      textinfo: 'label',
      insidetextorientation: 'radial'
    }];

    const layout: Partial<Plotly.Layout> = {
      title: 'Gráfico em Pizza com Legenda à Direita',
      height: 300,
      legend: { orientation: 'v', x: 0.62, y: 0.5, xanchor: 'left' },
      margin: { l: 50, r: 10, t: 50, b: 50 }
    };

    // se for re-render com os mesmos containers, prefira react
    if ((el as any).data) {
      Plotly.react(el, data, layout);
    } else {
      Plotly.newPlot(el, data, layout);
    }
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
           
                if(await this.invalidDate(this.searchAirlineDtInicio,  this.translate.instant('Ecra.invalidDateStart'), true)) {
                    return false;
                }
      
                if(await this.invalidDate(this.searchAirlineDtFim, this.translate.instant('Ecra.invalidDateEnd'), true)) {
                    return false;
                }
                    
                const startDate = moment(this.searchAirlineDtInicio, 'YYYY-MM-DD');
                const endDate = moment(this.searchAirlineDtFim, 'YYYY-MM-DD');
      
                if (startDate.isAfter(endDate)) {
                    const resultado = await this.modalOk.openModal(this.translate.instant('Ecra.StartDateGreaterThanEndDate'),true);             
                        if (resultado) {
                          return false;
                          // Insira aqui a lógica para continuar após a confirmação
                        } else {
                          
                        }
                  
                }
      
                if (endDate.isBefore(startDate)) {
                  const resultado = await this.modalOk.openModal(this.translate.instant('Ecra.StartDateEarlierThanEndDate'),true);             
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

        this.purcharseRecordService.GroupByCompanyName(objPesquisar.status,objPesquisar.companyCode,objPesquisar.startDate,objPesquisar.endDate).subscribe((response:any)=>{

         if(response) {
          console.log('adasd',response);
          response.forEach((element:any) => {
              _values.push(element.total)
              _labels.push(element._id + ' ' + element.total)
          });

          
          this.graphicData(_values,_labels,'graficoPlotlyDashBoard');
             
          this.msgNotFound=false;
        }
        else {
          //  console.log('not',response);
          this.msgNotFound=true;
          this.graphicData([],[],'graficoPlotlyDashBoard');
          
          
        }
      })

      return null
  }

  ngOnInit() {

      let _values:any[]=[]
      let _labels:any[]=[]

      this.purcharseRecordService.GroupByCompanyName().subscribe((response:any)=>{

       if(response) {
          response.forEach((element:any) => {
                  _values.push(element.total)
                  _labels.push(element._id + ' ' + element.total)
              });

              if(_values.length>0 && _labels.length>0)
                    this.graphicData(_values,_labels,'graficoPlotlyDashBoard');
            }
          
      })

        this.airLineService.getAllAirLine().subscribe((response)=>{        
          this.airLInes=this.configService.sortByKey(response,'name');      

        });

        this.moduloService.getAllStatusImportData().subscribe((response:any)=>{     
        this.statusImportData=response

          });
  }
}
