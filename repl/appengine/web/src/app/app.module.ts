import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule } from '@angular/common/http';
import { MatFormFieldModule } from '@angular/material/form-field'
import { MatInputModule } from '@angular/material/input'
import { MatSidenavModule } from '@angular/material/sidenav';
import {MatButtonModule} from '@angular/material/button';


import { AppComponent } from './app.component';
import { ReplConsoleComponent } from './repl-console/repl-console.component';
import { ReplResultDetailComponent } from './repl-result-detail/repl-result-detail.component';
import { TrimPipe } from './trim.pipe';
import { ReferencePanelComponent } from './reference-panel/reference-panel.component';

@NgModule({
  declarations: [
    AppComponent,
    ReplConsoleComponent,
    ReplResultDetailComponent,
    TrimPipe,
    ReferencePanelComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    MatFormFieldModule,
    BrowserAnimationsModule,
    MatInputModule,
    MatSidenavModule,
    MatButtonModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
