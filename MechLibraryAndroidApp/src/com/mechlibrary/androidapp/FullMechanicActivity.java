package com.mechlibrary.androidapp;

import com.koushikdutta.ion.Ion;

import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.ActionBarActivity;
import android.widget.ImageView;
import android.widget.TextView;

public class FullMechanicActivity extends ActionBarActivity {

	@Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.mechanic_full);
        Intent intent = getIntent();
        String title = intent.getStringExtra("title");
        String image = intent.getStringExtra("image");
        String description = intent.getStringExtra("description");
        String votes = intent.getStringExtra("votes");
        // start with the ImageView
        ImageView imageView = (ImageView)findViewById(R.id.image);
        Ion.with(imageView)
        // load the url
        .load(image);

        // and finally, set the name and text
        TextView handle = (TextView)findViewById(R.id.title);
        handle.setText(title + " - " + votes + " Votes");

        TextView text = (TextView)findViewById(R.id.description);
        text.setText(description);
    }

}
