
import { Component, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { ConfigService } from '../../../services/config.service';
import { ModalConfirmationComponent } from '../../../modal/modal-confirmation/modal-confirmation.component';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { PaginatorComponent } from '../../../paginator/paginator.component';

import { InvoicesService } from '../invoices/invoices.service';
import { ModalOkComponent } from '../../../modal/modal-ok/modal-ok.component';

@Component({
  selector: 'app-invoices',
  imports: [CommonModule,FormsModule,PaginatorComponent,ModalOkComponent],
  templateUrl: './invoices.component.html',
  styleUrl: './invoices.component.css'
})
export class InvoicesOnlyFlyComponent {

  

 
     @ViewChild(ModalOkComponent) modalOk!: ModalOkComponent;  


    searchKey: any = '';
    searchName: string = '';
    searchLastName: string = '';
    totalRegistros: number = 0;
    totalPages: number = 1;
    currentPage: number = 1;
    limit: number = 0;  
    currentYear: number = new Date().getFullYear();
    dados:any
    costCenterSubModalList:any[]=[]
    costCenters:any[]=[]    
    costCentersSub:any[]=[]

    constructor(
      private router: Router, 
      private invoiceService: InvoicesService,
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
            this.totalRegistros = response.total;
            this.totalPages = response.pages;
        });
     }
  
    
  
    onPageChange(newPage: number) {
      this.currentPage = newPage;
      this.searchOnlyFlyExcelData(this.currentPage);
    }
  
    
    searchOnlyFlyExcelData(currentPage:number) {    
  
      let objPesquisar: { 
          key: string;
          name: string;   
          lastName: string;            
          page:number;
          limit:number;
      }
  
      objPesquisar= { 
        key: this.searchKey, 
        name: this.searchName, 
        lastName: this.searchLastName,      
        page:currentPage,
        limit:this.limit
      };
  
      this.invoiceService.searchOnlyFlyExcelData(objPesquisar).subscribe((response:any)=>{
        this.dados=response.onlyFlyData;    
        this.totalRegistros = response.total;
        this.totalPages = response.pages;
      })
  
    }
  
    addInvoicesOnlyFlyExcel() {
      this.router.navigate(['/aplicacao/addInvoicesOnlyFlyExcel']);
    }
  
    updateInvoicesOnlyFlyExcel(id:string) {
      this.router.navigate(['/aplicacao/addInvoicesOnlyFlyExcel', id]);   
     } 
    
  

}
