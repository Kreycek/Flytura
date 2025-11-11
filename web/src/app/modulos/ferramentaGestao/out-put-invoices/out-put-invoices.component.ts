
import { CommonModule } from '@angular/common';
import { Component, ViewChild } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { TranslateModule } from '@ngx-translate/core';
import { AirLineService } from '../../companys/airLine/airLIne.service';
import { ConfigService } from '../../../services/config.service';
import moment from 'moment';
import { ModalOkComponent } from '../../../modal/modal-ok/modal-ok.component';
import { MAT_DATE_LOCALE, MatNativeDateModule } from '@angular/material/core';

import { PaginatorComponent } from '../../../paginator/paginator.component';
import { OutPutInvoicesService } from '../out-put-invoices/outPutInvoices.service';




@Component({
  selector: 'app-out-put-invoices',
  imports: [CommonModule,FormsModule,MatDatepickerModule,TranslateModule,MatNativeDateModule,ModalOkComponent,PaginatorComponent],
  templateUrl: './out-put-invoices.component.html',
  styleUrl: './out-put-invoices.component.css',
   providers: [{ provide: MAT_DATE_LOCALE, useValue: 'pt-BR' }]
})

export class OutPutInvoicesComponent {

  @ViewChild(ModalOkComponent) modalOk!: ModalOkComponent;  
   searchKeyCode:string='';
    searchDtInicio:string='';
    searchDtFim:string='';
    searchCompanyCode='';
   
    outPutInvoices:any[]=[];
    airLines:any[]=[];
    msgNotFound=false;
    totalRegistros: number = 0;
    totalPages: number = 1;    
    currentPage: number = 1;
    limit: number = 0;  
    dados:any;
   
    objPesquisar:any= { 
          
         
      };
    
    /**
     *
     */
    constructor(
      private airLineService: AirLineService,
            public configService:ConfigService,
            public outPutService:OutPutInvoicesService
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

    ngOnInit() {

        this.limit=this.configService.limitPaginator;
        this.airLineService.getAllAirLine().subscribe((response)=>{
         this.airLines=response;
        });
        this.search(this.currentPage);           
    }

    async search(pageNumber:number) {       
          
        if(this.searchDtInicio || this.searchDtFim ) {
        
            if(await this.invalidDate(this.searchDtInicio, 'Data de início inválida, ou se a data de fim estiver preenchida é necessário preencher uma data de início.', true)) {
                return false;
            }
  
            if(await this.invalidDate(this.searchDtFim, 'Data de fim inválida, ou se a data de início estiver preenchida é necessário preencher uma data de fim.', true)) {
                return false;
            }
                
            const startDate = moment(this.searchDtInicio, 'YYYY-MM-DD');
            const endDate = moment(this.searchDtFim, 'YYYY-MM-DD');
  
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

       console.log('this.searchDtInicio  ',this.searchDtInicio==undefined );
      console.log('this.searchDtFim  ',this.searchDtFim==undefined );
      this.objPesquisar= {
            keyCode:this.searchKeyCode,
            companyCode:this.searchCompanyCode,
            startDate:this.searchDtInicio  ? moment(this.searchDtInicio).format('YYYY-MM-DDTHH:mm:ss[Z]') : '',                               
            endDate: this.searchDtFim ? moment(this.searchDtFim).set({ hour: 23, minute: 59, second: 59 }).format('YYYY-MM-DDTHH:mm:ss[Z]')  : ''     
      };

      this.outPutInvoices=[]
      this.outPutService.getAllOutPutInvoicesDataPagination(
                    pageNumber,
                    this.limit,
                    this.objPesquisar.keyCode,
                    this.objPesquisar.companyCode,
                    this.objPesquisar.startDate,
                    this.objPesquisar.endDate
              ).subscribe((response:any)=>{
               
                if(response.outPutInvoices) {

                  this.outPutInvoices=response.outPutInvoices;
                  console.log('response.outPutInvoices ',response.outPutInvoices)
                  
                } else {
                   this.dados=null;
                }

                  this.totalRegistros = response.total;
                  this.totalPages = response.pages;
              })

        return null;
    }

     
    onPageChange(newPage: number) {
      this.currentPage = newPage;
      this.search(this.currentPage);
    }
    
}
