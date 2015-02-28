package com.mechlibrary.androidapp;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;

import com.koushikdutta.ion.Ion;
import com.koushikdutta.ion.Response;
import com.google.gson.JsonObject;

import com.koushikdutta.async.future.FutureCallback;

import android.app.ProgressDialog;
import android.content.Intent;
import android.graphics.Bitmap;
import android.os.Bundle;
import android.os.Environment;
import android.support.v4.app.Fragment;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.ImageView;
import android.widget.TextView;
import android.widget.Toast;

public class SubmitMechanic extends Fragment {
	public ImageView imgView;
	public Bitmap bitmap;
	public ProgressDialog dialog;
	final String server = "http://mech-library.appspot.com";
	//final String server = "10.0.2.2:8080";

    long totalSize = 0;

    public SubmitMechanic() {
    }
    
    public void setImage() {
    	if(bitmap!=null) {
    		imgView = (ImageView) getView().findViewById(R.id.upload_image);
			imgView.setImageBitmap(bitmap);
    	}
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
            Bundle savedInstanceState) {
    	View rootView = inflater.inflate(R.layout.submit_mechanic, container, false);
    	final TextView title = (TextView) rootView.findViewById(R.id.editText1);
    	final TextView description = (TextView) rootView.findViewById(R.id.editText2);
    	Button selectButton = (Button) rootView.findViewById(R.id.upload_button);
    	setImage();
		selectButton.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
            	Intent intent = new Intent();
            	intent.setType("image/*");
            	intent.setAction(Intent.ACTION_GET_CONTENT);
            	intent.putExtra("return-data", true); //added snippet
            	getActivity().startActivityForResult(Intent.createChooser(intent, "Select Picture"),100);
            }
        });
       	Button button = (Button) rootView.findViewById(R.id.submit_button);
        button.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
                // validate
           	 	if (bitmap == null) {
                    Toast.makeText(getActivity().getApplicationContext(),
                            "Please select image", Toast.LENGTH_SHORT).show();
	            } else {
	           	 // submit
	                dialog = ProgressDialog.show((LibraryActivity)getActivity(), "Uploading",
	                        "Please wait...", true);
	
	                final File file = new File(Environment.getExternalStorageDirectory() + "/myimage.png");
	                FileOutputStream fOut;
					try {
						fOut = new FileOutputStream(file);
	                    bitmap.compress(Bitmap.CompressFormat.PNG, 75, fOut);
	                    fOut.close();
					} catch (FileNotFoundException e) {
						// TODO Auto-generated catch block
						e.printStackTrace();
					} catch (IOException e) {
						// TODO Auto-generated catch block
						e.printStackTrace();
					}
					Ion.getDefault(getActivity().getApplicationContext()).configure().setLogging("MyLogs", Log.DEBUG);
					Ion.with(getActivity().getApplicationContext())
	            	.load(server+"/android/mechanics/add")
	            	.followRedirect(true)
	            	.asString()
	            	.withResponse()
	            	.setCallback(new FutureCallback<Response<String>>() {
	            	    @Override
	            	    public void onCompleted(Exception e, Response<String> result) {
	            	        // print the response code, ie, 200
	            	        System.out.println(result.getHeaders().code());
	            	        // print the String that was downloaded
	            	        System.out.println(result.getResult());
	            	        Ion.with(getActivity().getApplicationContext())
	    	            	.load(server+result.getResult())
	    	            	.followRedirect(true)
	    	            	.progressDialog(dialog)
	    	            	.setMultipartParameter("title", title.getText().toString())
	    	            	.setMultipartParameter("description", description.getText().toString())
	    	            	.setMultipartFile("screenshot", file)
	    	            	.asJsonObject()
	    	            	.withResponse()
	    	            	.setCallback(new FutureCallback<Response<JsonObject>>() {
	    	            	    @Override
	    	            	    public void onCompleted(Exception e, Response<JsonObject> result) {
	    	            	    	Toast.makeText(getActivity(), "Game Mechanic Uploaded! Thanks.", Toast.LENGTH_LONG).show();
	    	            	        dialog.dismiss();
	    	            	    }
	    	            	});
	            	    }
	            	});
	                
	            }
            }
        });
        return rootView;
    }
 
}
