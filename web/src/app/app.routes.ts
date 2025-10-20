import { RouterModule, Routes } from '@angular/router';
import { CenterComponent } from './center/center.component';
import { NgModule } from '@angular/core';

import { AplicacaoComponent } from './aplicacao/aplicacao.component';
import { LoginComponent } from './login/login.component';
import { UsuarioComponent } from './modulos/ferramentaGestao/users/usuario/usuario.component';
import { AddUsuarioComponent } from './modulos/ferramentaGestao/users/add-usuario/add-usuario.component';


import { PurcharseRecordComponent } from './modulos/companys/purcharseRecord/purcharse-record.component';
import { AddPurchaseRecordComponent } from './modulos/companys/add-purchase-record/add-purcharse-record.component';
import { ModelsComponent } from './modulos/companys/models/models/models.component';
import { InvoicesComponent } from './modulos/ferramentaGestao/invoices/invoices.component';

export const routes: Routes = [

    { path: '', redirectTo: '/login', pathMatch: 'full' },   // Redireciona para login por padrão    
    { 
        path:'aplicacao',
        component:AplicacaoComponent,
        children:[
          
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
                path:'addPurchaseRecord',
                component:AddPurchaseRecordComponent
            }, 
            {
                path:'addPurchaseRecord/:id',
                component:AddPurchaseRecordComponent
            }, 
            {
                path:'sheetModels',
                component:ModelsComponent
            } , 
            {
                path:'invoices',
                component:InvoicesComponent
            }

            
        ]
    },   // Redireciona para login por padrão    
    { 
        path:'login',
        component:LoginComponent    
    }
   
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
  })
  export class AppRoutingModule { }