import { RouterModule, Routes } from '@angular/router';
import { CenterComponent } from './center/center.component';
import { NgModule } from '@angular/core';

import { AplicacaoComponent } from './aplicacao/aplicacao.component';
import { LoginComponent } from './login/login.component';
import { UsuarioComponent } from './modulos/gestion/users/list/usuario.component';
import { AddUsuarioComponent } from './modulos/gestion/users/addUpdate/add-usuario.component';


import { PurcharseRecordComponent } from './modulos/Account/purcharseRecord/list/purcharse-record.component';
import { AddPurchaseRecordComponent } from './modulos/Account/purcharseRecord/addUpdate/add-purcharse-record.component';
import { ModelsComponent } from './modulos/gestion/sheetModels/models.component';
import { InvoicesComponent } from './modulos/gestion/invoices/list/invoices.component';
import { OutPutInvoicesComponent } from './modulos/gestion/out-put-invoices/out-put-invoices.component';
import { HashLocationStrategy, LocationStrategy } from '@angular/common';
import { SectionComponent } from './modulos/gestion/section/section.component';
import { ConciliationListComponent } from './modulos/Account/conciliation/list/list.component';

export const routes: Routes = [

    { path: '', redirectTo: 'login', pathMatch: 'full' },   // Redireciona para login por padrão    
    { 
        path:'aplicacao',
        component:AplicacaoComponent,
        children:[
          
            {
                path:'section',
                component:SectionComponent,
            },
            {
                path:'center',
                component:CenterComponent,
            },
            {
                path:'usuario',
                component:UsuarioComponent,
            }, 
            {
                path:'addUser/:id',
                component:AddUsuarioComponent
            }   , 
            {
                path:'addUser',
                component:AddUsuarioComponent
            },             
            {
                path:'purcharseRecord',
                component:PurcharseRecordComponent
            },             
            {
                path:'purcharseRecord/:key/:nombre/:apellido/:companyCodAirline/:dtStart/:dtEnd/:status',
                component:PurcharseRecordComponent
            }, 
            {
                path:'outPutInvoices',
                component:OutPutInvoicesComponent
            },            
            {
                path:'addPurchaseRecord',
                component:AddPurchaseRecordComponent
            }, 
            {
                path:'addPurchaseRecord/:id/:key/:nombre/:apellido/:companyCodAirline/:dtStart/:dtEnd/:status',
                component:AddPurchaseRecordComponent
            }, 
            {
                path:'sheetModels',
                component:ModelsComponent
            } , 
            {
                path:'invoices',
                component:InvoicesComponent
            } , 
            {
                path:'conciliation',
                component:ConciliationListComponent
            }

            
        ]
    },   // Redireciona para login por padrão    
    { 
        path:'login/:codAKD',
        component:LoginComponent    
    },
    { 
        path:'login',
        component:LoginComponent    
    },   // Redireciona para login por padrão    
    
   
];

@NgModule({
    imports: [RouterModule.forRoot(routes, { useHash: true })],
    exports: [RouterModule],
    providers: [{ provide: LocationStrategy, useClass: HashLocationStrategy }]

  })
  export class AppRoutingModule { }