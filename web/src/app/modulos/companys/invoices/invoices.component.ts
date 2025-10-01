
import { Component, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { ConfigService } from '../../../services/config.service';
import { ModalConfirmationComponent } from '../../../modal/modal-confirmation/modal-confirmation.component';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { PaginatorComponent } from '../../../paginator/paginator.component';
import * as _moment from 'moment';
import { InvoicesService } from '../invoices/invoices.service';
import { ModalOkComponent } from '../../../modal/modal-ok/modal-ok.component';
import { AirLineService } from '../airLine/airLIne.service';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MAT_DATE_LOCALE, MatNativeDateModule } from '@angular/material/core';
import moment from 'moment';

@Component({
  selector: 'app-invoices',
  imports: [CommonModule,FormsModule,PaginatorComponent,ModalOkComponent,MatDatepickerModule,MatNativeDateModule,],
  templateUrl: './invoices.component.html',
  styleUrl: './invoices.component.css',
  providers: [{ provide: MAT_DATE_LOCALE, useValue: 'pt-BR' }],
})
export class InvoicesOnlyFlyComponent { 
     @ViewChild(ModalOkComponent) modalOk!: ModalOkComponent;  


    searchKey: any = '';
    searchName: string = '';
    searchLastName: string = '';
    searchAirlineCode:string='';
    searchAirlineDtInicio:string='';
    searchAirlineDtFim :string='';
     searchStatus:string='';
    totalRegistros: number = 0;
    totalPages: number = 1;
    currentPage: number = 1;
    limit: number = 0;  
    currentYear: number = new Date().getFullYear();
    dados:any
     statusImportData:any[]=[]
    costCenterSubModalList:any[]=[]
    costCenters:any[]=[]    
    costCentersSub:any[]=[]
 airLInes:any[]=[]

    constructor(
      private router: Router, 
      private invoiceService: InvoicesService,
      private airLineService: AirLineService,
      public configService:ConfigService
      
    ) {} 
      
  
      
   onFileSelected(event: any) {
      const file: File = event.target.files[0];
     
      if (file) {
        const formData = new FormData();
              formData.append('file', file);
              // console.log(file.name);
            this.invoiceService.importarPlanilha(formData).subscribe(()=>{
                this.invoiceService.getAllOnlyFlyExcelDataPagination(this.currentPage,this.limit).subscribe(async (response:any)=>{     
        
                  this.dados=response.onlyFlyData;  
                  this.totalRegistros = response.total;
                  this.totalPages = response.pages;
                    const resultado = await this.modalOk.openModal("Planilha importada com sucesso",true);             
                      if (resultado) {
                    
                        // Insira aqui a lógica para continuar após a confirmação
                      } else {
                        
                      }
              });
            })
  
          }
  
    }
  
   
      ngOnInit() {
   
        this.limit=this.configService.limitPaginator;
     
        this.invoiceService.getAllOnlyFlyExcelDataPagination(this.currentPage,this.limit).subscribe((response:any)=>{     
  
            this.dados=response.onlyFlyData;  
            console.log('this.dados ',this.dados);
            this.totalRegistros = response.total;
            this.totalPages = response.pages;
        });

        this.airLineService.getAllAirLine().subscribe((response)=>{
          this.airLInes=response;
          console.log('teste',response);
        });

this.invoiceService.getAllStatusImportData().subscribe((response:any)=>{     
        this.statusImportData=response

          });
     }
  
    
  
    onPageChange(newPage: number) {
      this.currentPage = newPage;
      this.searchOnlyFlyExcelData(this.currentPage);
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
    
    async searchOnlyFlyExcelData(currentPage:number) {    
  
      let objPesquisar: { 
          key: string;
          name: string;   
          lastName: string;    
          companyCode:string;  
          startDate?:string | null;
          endDate?:string | null;
          status?:string | null;
          page:number;
          limit:number;
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
        key: this.searchKey, 
        name: this.searchName, 
        lastName: this.searchLastName, 
        companyCode:this.searchAirlineCode,
        startDate:this.searchAirlineDtInicio ? moment(new Date(this.searchAirlineDtInicio)).format('YYYY-MM-DDT00:00:00Z') : null,
        endDate:this.searchAirlineDtFim ? moment(new Date(this.searchAirlineDtFim)).format('YYYY-MM-DDT23:59:59Z') :null,
        status:this.searchStatus,
        page:currentPage,
        limit:this.limit
      };
      
  
      this.invoiceService.searchOnlyFlyExcelData(objPesquisar).subscribe((response:any)=>{
        this.dados=response.onlyFlyData;    
        this.totalRegistros = response.total;
        this.totalPages = response.pages;
      })

      return null;
  
    }
  
    addInvoicesOnlyFlyExcel() {
      this.router.navigate(['/aplicacao/addInvoicesOnlyFlyExcel']);
    }
  
    updateInvoicesOnlyFlyExcel(id:string) {
      this.router.navigate(['/aplicacao/addInvoicesOnlyFlyExcel', id]);   
     } 
    
  
  
tooltip = {
    visible: false,
    text: '',
    top: 0,
    left: 0,
    maxWidth: 'auto'
  };

  mousemoveHandler(event: MouseEvent, item: any, show:boolean=true, textMessage:string) {
    const offsetX = 15;
    const offsetY = 15;
    const tooltipWidth = 200;
    const screenWidth = window.innerWidth;

    let left = event.clientX + offsetX;
    let top = event.clientY + offsetY;
    let maxWidth = 'auto';

    if (left + tooltipWidth > screenWidth) {
      left = screenWidth - tooltipWidth - offsetX;
      maxWidth = '180px';
    }

    this.tooltip = {
      visible: show,
      text:textMessage,
      top,
      left,
      maxWidth
    };
  }

  hideCustomTooltip() {
    this.tooltip.visible = false;
  }


}
