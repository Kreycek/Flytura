
import { CommonModule } from '@angular/common';
import { Component, ElementRef, LOCALE_ID, ViewChild } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { AirLineService } from '../airLine/airLIne.service';
import { ConfigService } from '../../../services/config.service';
import moment from 'moment';
import { ModalOkComponent } from '../../../modal/modal-ok/modal-ok.component';
import { MAT_DATE_LOCALE, MatNativeDateModule } from '@angular/material/core';

import { PaginatorComponent } from '../../../paginator/paginator.component';
import { OutPutInvoicesService } from './outPutInvoices.service';


import { registerLocaleData } from '@angular/common';
import localePt from '@angular/common/locales/pt';
import { ModalConfirmationComponent } from '../../../modal/modal-confirmation/modal-confirmation.component';
import { jwtDecode } from 'jwt-decode';
import { EnumPerfil } from '../../Enum/perfil';
import { AlertMoreColumnsComponent } from "../../../components/alert-more-columns/alert-more-columns.component";
import { ModuloService } from '../../modulo.service';
import { SpinnerComponent } from '../../../components/spinner/spinner.component';
registerLocaleData(localePt, 'pt');



@Component({
  selector: 'app-out-put-invoices',
  imports: [
    CommonModule, 
    FormsModule, 
    MatDatepickerModule, 
    TranslateModule, 
    MatNativeDateModule, 
    ModalOkComponent, 
    PaginatorComponent, 
    TranslateModule, 
    ModalConfirmationComponent, 
    AlertMoreColumnsComponent,
    SpinnerComponent
  ],
  templateUrl: './out-put-invoices.component.html',
  styleUrl: './out-put-invoices.component.css',
   providers: [{ provide: MAT_DATE_LOCALE, useValue: 'pt-BR' },{ provide: LOCALE_ID, useValue: 'pt' }]
})

export class OutPutInvoicesComponent {
 @ViewChild('fileInput') fileInput!: ElementRef;
  @ViewChild(ModalOkComponent) modalOk!: ModalOkComponent; 
  
     @ViewChild(ModalConfirmationComponent) modalConfirm!: ModalConfirmationComponent; 
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
        objPesquisar:any= { };
        adm:boolean=true;
        fillOneFilter=false;
        isLoading:boolean=false;
        decoded:any
    
    /**
     *
     */
    constructor(
      private airLineService: AirLineService,
            public configService:ConfigService,
            public outPutService:OutPutInvoicesService,
                  private translate: TranslateService,
                  public moduloService:ModuloService
          ) {}

          
    verifyFieldsSearchFill() {   
      
      const campos = [
        this.searchKeyCode,
        this.searchDtInicio,
        this.searchDtFim,
        this.searchCompanyCode,
      
      ];

      this.fillOneFilter = campos.some(v => v !== null && v != '');

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
            this.decoded = jwtDecode<any>(token);
            console.log('decode ',this.decoded);
            if(this.decoded && this.decoded.perfis && this.decoded.perfis.length>0) {
                this.adm=this.decoded.perfis.some((data:any)=> data===EnumPerfil.ADM)
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

    
    async deleteItem(id:string) {

        const resultado = await this.modalConfirm.openModal(
                  true,
                  this.translate.instant('Ecra.deleteInformation'),
                  this.translate.instant('Ecra.yes'),
                  this.translate.instant('Ecra.no')
                ); 
        if (resultado) {     
          
              this.modalConfirm.isVisible=false;
              this.outPutService.deleteOutPutInvoices(id).subscribe();
              this.search(this.currentPage);        
        } else {
          this.modalConfirm.isVisible=false;
      
        }
    }
    
    
        generateReportExcel() {    

             this.objPesquisar= {
                keyCode:this.searchKeyCode,
                companyCode:this.searchCompanyCode,
                startDate:this.searchDtInicio  ? moment(this.searchDtInicio).format('YYYY-MM-DDTHH:mm:ss[Z]') : '',                               
                endDate: this.searchDtFim ? moment(this.searchDtFim).set({ hour: 23, minute: 59, second: 59 }).format('YYYY-MM-DDTHH:mm:ss[Z]')  : ''     
            };

            this.outPutInvoices=[]
            this.outPutService.getAllOutPutInvoicesDataInforme(                  
                          this.objPesquisar.keyCode,
                          this.objPesquisar.companyCode,
                          this.objPesquisar.startDate,
                          this.objPesquisar.endDate
                    ).subscribe((response:any)=>{               
                        if(response.outPutInvoices) {
                          this.outPutInvoices=response.outPutInvoices;  
                          this.moduloService.exportToExcelInformes(this.outPutInvoices,'informes')
                        
                        } else {
                          this.dados=null;
                        }                 
                    })


        
        }


        
  
  onFileSelected(event: any) {
      const file: File = event.target.files[0];
     
      if (file) {
              const formData = new FormData();
              formData.append('file', file);
              formData.append('idUserInserted', this.decoded.idUser);
              formData.append('userName', this.decoded.name);
              this.isLoading=true;

          this.moduloService.importFile(formData, "/UploadExcelOutputInvoicesRecord").subscribe({
            next: async (returnSheet: any) => {
              if (returnSheet.message) {
                console.log('returnSheet.message ',returnSheet.message);
                // Caso 1: Planilha vazia
                  this.isLoading=false;
                  this.fileInput.nativeElement.value = '';
                   const resultado = await this.modalOk.openModal(this.translate.instant('Ecra.sheetImportedOk'), true);
                  if (resultado) {
                    this.fileInput.nativeElement.value = '';
                  }
                 this.search(1); 
              }
            },
            error: async (err: any) => {
              console.error('Erro do servidor:', err);
              if (err.error?.status === "500") {
                 this.isLoading=false;
                const resultado = await this.modalOk.openModal(                
                   this.translate.instant('Ecra.correctSheetName')
                   .replace("{0}", err.error?.codFile)
                   .replace("{1}", err.error?.sheetName)
                   .replace("{2}", err.error?.codFile)
                   .replace("{3}", err.error?.nameAirLine)
                   .replace("{4}", err.error?.codFile)
                    .replace("{5}", err.error?.fileName)
                  , true);
                if (resultado) {
                    
                  this.fileInput.nativeElement.value = '';
                }
              }
              else if (err.status === 500) {
                   this.isLoading=false;
                const resultado = await this.modalOk.openModal(                
                    this.translate.instant(err.error.replace("\n","")), 
                    true);
                     if (resultado) {
                      
                  this.fileInput.nativeElement.value = '';
                }
              }
            }
          });
        } 
    }

downloadFile() {
  const link = document.createElement('a');
  link.href = '../../../../../assets/airLineSheetModel/SheetOutputInvoicesModel.xlsx';
  link.download = 'SheetOutputInvoicesModel.xlsx';
  link.click();
}
   
    
}
