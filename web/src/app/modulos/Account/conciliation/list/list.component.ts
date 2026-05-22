import { CommonModule } from '@angular/common';
import { Component, ElementRef, LOCALE_ID, ViewChild } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { AirLineService } from '../../../gestion/airLine/airLIne.service';
import { ConfigService } from '../../../../services/config.service';
import moment from 'moment';
import { ModalOkComponent } from '../../../../modal/modal-ok/modal-ok.component';
import { MAT_DATE_LOCALE, MatNativeDateModule } from '@angular/material/core';

import { PaginatorComponent } from '../../../../paginator/paginator.component';
import { ConciliationService } from '../conciliation.service';


import { registerLocaleData } from '@angular/common';
import localePt from '@angular/common/locales/pt';
import { ModalConfirmationComponent } from '../../../../modal/modal-confirmation/modal-confirmation.component';
import { jwtDecode } from 'jwt-decode';
import { EnumPerfil } from '../../../Enum/perfil';
import { AlertMoreColumnsComponent } from "../../../../components/alert-more-columns/alert-more-columns.component";
import { ModuloService } from '../../../modulo.service';
registerLocaleData(localePt, 'pt');
@Component({
  selector: 'app-list',
  imports: [CommonModule, FormsModule, MatDatepickerModule, TranslateModule, MatNativeDateModule, ModalOkComponent, PaginatorComponent, TranslateModule, ModalConfirmationComponent, AlertMoreColumnsComponent],
  templateUrl: './list.component.html',
  styleUrl: './list.component.css'
})
export class ConciliationListComponent {

  @ViewChild(ModalOkComponent) modalOk!: ModalOkComponent; 
  
     @ViewChild(ModalConfirmationComponent) modalConfirm!: ModalConfirmationComponent; 
    searchOriginLocator:string='';
    searchReturnLocator:string='';
    searchOriginETicket:string='';
    searchReturnETicket:string='';
    searchDtInicio:string='';
    searchDtFim:string='';
    searchCompanyCode='';   
    conciliation:any[]=[];
    airLines:any[]=[];
    msgNotFound=false;
    totalRegistros: number = 0;
    totalPages: number = 1;    
    currentPage: number = 1;
    limit: number = 0;  
    dados:any;   
    objPesquisar:any= { };
    adm:boolean=true;
      fillOneFilter=false;
    
    /**
     *
     */
    constructor(
      private airLineService: AirLineService,
            public configService:ConfigService,
            public conciliartionService:ConciliationService,
                  private translate: TranslateService,
                  public moduloService:ModuloService
          ) {}

          
    verifyFieldsSearchFill() {

   
      
      const campos = [
        this.searchOriginLocator,
        this.searchReturnLocator,
        this.searchOriginETicket,
        this.searchReturnETicket,
        this.searchDtInicio,
        this.searchDtFim,
        this.searchCompanyCode,
      
      ];

      this.fillOneFilter = campos.some(v => v !== null && v != '');

      console.log('this.existeAlgumPreenchido',this.fillOneFilter);

    }

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
      
      const token = localStorage.getItem('token'); // ou onde você armazenou o JWT
      
          if (token) {
            const decoded = jwtDecode<any>(token);
            if(decoded && decoded.perfis && decoded.perfis.length>0) {
                this.adm=decoded.perfis.some((data:any)=> data===EnumPerfil.ADM)
            }
          }


        this.limit=this.configService.limitPaginator;
        this.airLineService.getAllAirLine().subscribe((response)=>{
           this.airLines=this.configService.sortByKey(response,'name');   
       
         
        });
        this.search(this.currentPage);           

        // this.outPutService.UpdateWay({key:'CCK6ZH',way:120}).subscribe()
    }

    async search(pageNumber:number) {       
          
        if(this.searchDtInicio || this.searchDtFim ) {
        
            if(await this.invalidDate(this.searchDtInicio, this.translate.instant('Ecra.invalidDateStart'), true)) {
                return false;
            }
  
            if(await this.invalidDate(this.searchDtFim,  this.translate.instant('Ecra.invalidDateEnd'), true)) {
                return false;
            }
                
            const startDate = moment(this.searchDtInicio, 'YYYY-MM-DD');
            const endDate = moment(this.searchDtFim, 'YYYY-MM-DD');
  
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

      //  console.log('this.searchDtInicio  ',this.searchDtInicio==undefined );
      // console.log('this.searchDtFim  ',this.searchDtFim==undefined );
      this.objPesquisar= {
            originLocator:this.searchOriginLocator,
            returnLocator:this.searchReturnLocator,
            originETicket:this.searchOriginETicket,
            returnETicket:this.searchReturnETicket,
            startDate:this.searchDtInicio  ? moment(this.searchDtInicio).format('YYYY-MM-DDTHH:mm:ss[Z]') : '',                               
            endDate: this.searchDtFim ? moment(this.searchDtFim).set({ hour: 23, minute: 59, second: 59 }).format('YYYY-MM-DDTHH:mm:ss[Z]')  : ''     
      };

      this.conciliation=[]
      this.conciliartionService.getAllConciliationDataPagination(
                    pageNumber,
                    this.limit,
                    this.objPesquisar.originLocator,
                    this.objPesquisar.returnLocator,
                    this.objPesquisar.originETicket,
                    this.objPesquisar.returnETicket,                 
                    this.objPesquisar.startDate,
                    this.objPesquisar.endDate
              ).subscribe((response:any)=>{               
                  if(response.conciliation) {
                    this.conciliation=response.conciliation;  
                   
                  } else {
                    this.dados=null;
                  }

                  this.totalRegistros = response.total;
                  if(this.totalRegistros==0) {
                      this.currentPage=1;
                  }
                  this.totalPages = response.pages;
              })


              

        return null;
    }

     
    onPageChange(newPage: number) {
      this.currentPage = newPage;
      this.search(this.currentPage);
    }   
    
    
       async  generateReportExcel() {    


           if(this.searchDtInicio || this.searchDtFim ) {
        
            if(await this.invalidDate(this.searchDtInicio, this.translate.instant('Ecra.invalidDateStart'), true)) {
                return false;
            }
  
            if(await this.invalidDate(this.searchDtFim,  this.translate.instant('Ecra.invalidDateEnd'), true)) {
                return false;
            }
                
            const startDate = moment(this.searchDtInicio, 'YYYY-MM-DD');
            const endDate = moment(this.searchDtFim, 'YYYY-MM-DD');
  
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

             this.objPesquisar= {
                originLocator:this.searchOriginLocator,
                returnLocator:this.searchReturnLocator,
                originETicket:this.searchOriginETicket,
                returnETicket:this.searchReturnETicket,
                startDate:this.searchDtInicio  ? moment(this.searchDtInicio).format('YYYY-MM-DDTHH:mm:ss[Z]') : '',                               
                endDate: this.searchDtFim ? moment(this.searchDtFim).set({ hour: 23, minute: 59, second: 59 }).format('YYYY-MM-DDTHH:mm:ss[Z]')  : ''     
            };

            // this.conciliation=[]
            this.conciliartionService.getAllConciliationExcel(                  
                          this.objPesquisar.originLocator,
                          this.objPesquisar.returnLocator,
                          this.objPesquisar.originETicket,
                          this.objPesquisar.returnETicket,
                          this.objPesquisar.startDate,
                          this.objPesquisar.endDate
                    ).subscribe((response:any)=>{     
                      console.log('response.conciliation',response.conciliation);          
                        if(response.conciliation) {
                          // this.conciliation=response.outPutInvoices;  
                          this.moduloService.exportToExcelConciliation(response.conciliation,'conciliation')
                        
                        } else {
                          this.dados=null;
                        }                 
                    })

                    return null
        
        }
}
