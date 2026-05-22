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
  selector: 'app-grafic-import-and-status',
  imports: [CommonModule, FormsModule, MatDatepickerModule, TranslateModule, MatNativeDateModule, ModalOkComponent],
   providers: [{ provide: MAT_DATE_LOCALE, useValue: 'pt-BR' }],
  templateUrl: './grafic-import-and-status.component.html',
  styleUrl: './grafic-import-and-status.component.css'
})
export class GraficImportAndStatusComponent {
@ViewChild(ModalOkComponent) modalOk!: ModalOkComponent;  

    searchAirlineCode:string='';
    searchAirlineDtInicio:string='';
    searchAirlineDtFim :string='';
    searchStatus:string='';
    statusImportData:any[]=[];
    airLInes:any[]=[];
    msgNotFound=false;
    data:any[]=[];

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
generateLineChartFromList(list: any[]) {

  const grouped: Record<string, { x: string[]; y: number[]; total: number }> = {};

  // 1️⃣ Agrupar por status
  list.forEach(item => {
    if (!grouped[item.status]) {
      grouped[item.status] = { x: [], y: [], total: 0 };
    }
    grouped[item.status].x.push(item.date);
    grouped[item.status].y.push(item.total);
    grouped[item.status].total += item.total; // ✅ soma por status
  });

  // 2️⃣ Criar uma linha por status
  const traces = Object.keys(grouped).map(status => ({
    x: grouped[status].x,
    y: grouped[status].y,
    type: 'scatter',
    mode: 'lines+markers',

    // ✅ LEGENDA COM VALOR (FORA DO GRÁFICO)
    name: `${status} (${grouped[status].total})`,

    marker: { size: 6 },     // bolinhas normais
    line: { width: 3 },

    hovertemplate:
      'Data: %{x|%d/%m/%Y}<br>' +
      'Total: %{y}<extra></extra>'
  }));

  // 3️⃣ Layout com legenda fora
  const layout = {
    title: 'Evolução por Data e Status',
    xaxis: {
      title: 'Data',
      type: 'date',
      tickformat: '%d/%m/%Y'
    },
    yaxis: {
      title: 'Total'
    },
    legend: {
      orientation: 'v',
      x: 1.05,   // ✅ fora do gráfico (direita)
      y: 1
    }
  };

  Plotly.react('graficoLinha', traces, layout);
}


generateAreaChartFromList(list: any[]) {

  const grouped: Record<string, { x: string[]; y: number[]; total: number }> = {};

  // Agrupar por status
  list.forEach(item => {
    if (!grouped[item.status]) {
      grouped[item.status] = { x: [], y: [], total: 0 };
    }
    grouped[item.status].x.push(item.date);
    grouped[item.status].y.push(item.total);
    grouped[item.status].total += item.total;
  });

  const traces = Object.keys(grouped).map(status => ({
    x: grouped[status].x,
    y: grouped[status].y,

    type: 'scatter',
    mode: 'lines',

    // ✅ ÁREA
    fill: 'tozeroy',

    // ✅ LINHA ARREDONDADA
    line: {
      shape: 'spline',   // 🔥 deixa arredondado
      smoothing: 1.2,    // 🔥 nível de suavização
      width: 2
    },

    // ✅ legenda fora com valor
    name: `${status} (${grouped[status].total})`,

    hovertemplate:
      'Data: %{x|%d/%m/%Y}<br>' +
      'Total: %{y}<extra></extra>'
  }));

  const layout = {
    title: 'Evolução por Data e Status',
    xaxis: {
      type: 'date',
      tickformat: '%d/%m/%Y'
    },
    yaxis: {
      title: 'Total'
    },
    legend: {
      x: 1.05,
      y: 1
    }
  };

  Plotly.react('graficoLinha', traces, layout);
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

            this.data=[]
             objPesquisar= { 
                  
                    companyCode:this.searchAirlineCode,
                    startDate:this.searchAirlineDtInicio ? moment(new Date(this.searchAirlineDtInicio)).format('YYYY-MM-DDT00:00:00Z') : null,
                    endDate:this.searchAirlineDtFim ? moment(new Date(this.searchAirlineDtFim)).format('YYYY-MM-DDT23:59:59Z') :null,
                    status:this.searchStatus
                  };


        this.purcharseRecordService.AgregateByImportDateAndStatus(objPesquisar.status,objPesquisar.companyCode,objPesquisar.startDate,objPesquisar.endDate).subscribe((resAgre:any[])=>{
         console.log('resAgre ',resAgre);

            let total=0;
            let status=''
            let date=''

            if(resAgre && resAgre.length>0) {
              for(let i=0;i<resAgre.length; i++) {
                          const fields=resAgre[i]

                          total=0;
                            status=''
                            date=''

                          for(let j=0;j<fields.length; j++) {
                            

                            if(fields[j].Key=='total') {
                              total=fields[j].Value
                            }

                            if(fields[j].Key=='date') {
                              date=fields[j].Value
                            }
                          
                            if(fields[j].Key=='status') {
                              status=fields[j].Value
                            }

                          }

                            this.data.push(
                                {
                                  total: total,
                                  date: date,
                                  status: status
                                }
                            )
                      

                          
                      }
                           this.msgNotFound=false;
                        console.log('this.data ',this.data);

                        this.generateAreaChartFromList(this.data)
                    }else {
          //  console.log('not',response);
          this.msgNotFound=true;
          // this.graphicData([],[],'graficoPlotlyDashBoard');
          
          
        }
        })

     

      return null
  }

  ngOnInit() {

      let _values:any[]=[]
      let _labels:any[]=[]

        const  objPesquisar= { 
                  
                    companyCode:this.searchAirlineCode,
                    startDate:this.searchAirlineDtInicio ? moment(new Date(this.searchAirlineDtInicio)).format('YYYY-MM-DDT00:00:00Z') : null,
                    endDate:this.searchAirlineDtFim ? moment(new Date(this.searchAirlineDtFim)).format('YYYY-MM-DDT23:59:59Z') :null,
                    status:this.searchStatus
                  };

       this.purcharseRecordService.AgregateByImportDateAndStatus(objPesquisar.status,objPesquisar.companyCode,objPesquisar.startDate,objPesquisar.endDate).subscribe((resAgre:any[])=>{
         console.log('resAgre ',resAgre);

            let total=0;
            let status=''
            let date=''

         for(let i=0;i<resAgre.length; i++) {
            const fields=resAgre[i]

             total=0;
               status=''
               date=''

             for(let j=0;j<fields.length; j++) {

               console.log('(fields[j] ',fields[j]);

              

              if(fields[j].Key=='total') {
                total=fields[j].Value
              }

              if(fields[j].Key=='date') {
                date=fields[j].Value
              }
             
               if(fields[j].Key=='status') {
                status=fields[j].Value
              }

             }

              console.log(total,date,status);

              this.data.push(
                  {
                    total: total,
                    date: date,
                    status: status
                  }
              )
         

            
         }

         console.log('this.data ',this.data);

         this.generateAreaChartFromList(this.data)
        })

         

     
        this.airLineService.getAllAirLine().subscribe((response)=>{        
          this.airLInes=this.configService.sortByKey(response,'name');      

        });

        this.moduloService.getAllStatusImportData().subscribe((response:any)=>{     
        this.statusImportData=response

          });
  }
}

